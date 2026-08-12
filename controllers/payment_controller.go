package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/services"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreatePaymentOrder(c *gin.Context) {

	orderID, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid order id"})
		return
	}
	userID := c.MustGet("userID").(uint)

	var order models.Order

	if err := config.DB.
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {

		c.JSON(404, gin.H{
			"error": "Order not found",
		})
		return
	}

	if order.PaymentStatus == "paid" || order.PaymentMethod == "cod" {
		c.JSON(400, gin.H{
			"error": "Order already paid",
		})
		return
	}

	amount := order.TotalAmount

	if order.RazorpayOrderID != "" {

		c.JSON(200, gin.H{
			"order_id":          order.ID,
			"razorpay_order_id": order.RazorpayOrderID,
			"amount":            amount,
			"currency":          "INR",
			"key":               os.Getenv("RAZORPAY_KEY_ID"),
		})
		return
	}

	client := config.NewRazorpayClient()

	data := map[string]interface{}{
		"amount":   int(amount * 100), // convert to paise
		"currency": "INR",
		"receipt":  fmt.Sprintf("order_%d", order.ID),
	}

	rzOrder, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("===== RAZORPAY ORDER CREATED =====")
	fmt.Println(rzOrder)

	rzOrderID, ok := rzOrder["id"].(string)
	if !ok {
		c.JSON(500, gin.H{
			"error": "Invalid Razorpay response",
		})
		return
	}

	order.RazorpayOrderID = rzOrderID
	order.PaymentStatus = "pending"

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save order",
		})
		return
	}

	c.JSON(200, gin.H{
		"order_id":          order.ID,
		"razorpay_order_id": order.RazorpayOrderID,
		"amount":            amount,
		"currency":          "INR",
		"key":               os.Getenv("RAZORPAY_KEY_ID"),
	})
	return

}

type VerifyRequest struct {
	OrderID   string `json:"order_id"`   //razorpay order id
	PaymentID string `json:"payment_id"` //razorpay payment id
	Signature string `json:"signature"`  //razorpay signature
}

func VerifyPayment(c *gin.Context) {

	var req VerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	message := req.OrderID + "|" + req.PaymentID

	h := hmac.New(sha256.New, []byte(os.Getenv("RAZORPAY_KEY_SECRET")))
	h.Write([]byte(message))

	generatedSignature := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal(
		[]byte(generatedSignature),
		[]byte(req.Signature),
	) {
		c.JSON(400, gin.H{"error": "Invalid signature"})
		return
	}

	var order models.Order
	userID := c.MustGet("userID").(uint)

	if err := config.DB.
		Where(
			"razorpay_order_id = ? AND user_id = ?",
			req.OrderID,
			userID,
		).
		First(&order).Error; err != nil {

		c.JSON(404, gin.H{
			"error": "Order not found",
		})
		return
	}

	if order.PaymentStatus == "paid" {
		c.JSON(400, gin.H{
			"error": "Order already paid",
		})
		return
	}

	order.PaymentStatus = "paid"
	order.Status = "processing"
	order.RazorpayPaymentID = req.PaymentID

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User //getting user details for sending email

	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "User not found",
		})
	}

	go services.SendOrderConfirmationEmail(
		user.Email,
		order.ID,
		order.TotalAmount,
	)

	c.JSON(200, gin.H{
		"message": "Payment verified successfully",
	})

}

func RazorpayWebhook(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid body",
		})
		return
	}

	fmt.Println("===== WEBHOOK RECEIVED =====")
	fmt.Println(string(body))

	signature := c.GetHeader("X-Razorpay-Signature")

	if signature == "" {
		c.JSON(400, gin.H{
			"error": "Missing webhook signature",
		})
		return
	}

	h := hmac.New(
		sha256.New,
		[]byte(os.Getenv("RAZORPAY_WEBHOOK_SECRET")),
	)

	h.Write(body)

	generatedSignature :=
		hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal(
		[]byte(generatedSignature),
		[]byte(signature),
	) {
		c.JSON(400, gin.H{
			"error": "Invalid webhook signature",
		})
		return
	}

	var payload map[string]interface{}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid payload",
		})
		return
	}
	event, ok := payload["event"].(string)
	if !ok {
		c.JSON(400, gin.H{
			"error": "Event not found",
		})
		return
	}

	if event != "payment.captured" {

		c.JSON(200, gin.H{
			"message": "Event ignored",
		})

		return
	}

	payloadData, ok := payload["payload"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{
			"error": "Invalid payload structure",
		})
		return
	}

	paymentData, ok := payloadData["payment"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{
			"error": "Payment data missing",
		})
		return
	}

	entityData, ok := paymentData["entity"].(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{
			"error": "Payment entity missing",
		})
		return
	}
	razorpayOrderID, ok := entityData["order_id"].(string)
	if !ok {
		c.JSON(400, gin.H{
			"error": "Order ID missing",
		})
		return
	}

	razorpayPaymentID, ok := entityData["id"].(string)
	if !ok {
		c.JSON(400, gin.H{
			"error": "Payment ID missing",
		})
		return
	}

	var order models.Order

	if err := config.DB.
		Where("razorpay_order_id = ?", razorpayOrderID).
		First(&order).Error; err != nil {

		c.JSON(404, gin.H{
			"error": "Order not found",
		})
		return
	}

	if order.PaymentStatus == "paid" {

		c.JSON(200, gin.H{
			"message": "Already processed",
		})

		return
	}
	order.RazorpayPaymentID = razorpayPaymentID
	order.PaymentStatus = "paid"
	order.Status = "processing"
	if err := config.DB.Save(&order).Error; err != nil {

		c.JSON(500, gin.H{
			"error": "Failed to update order",
		})

		return
	}
	c.JSON(200, gin.H{
		"message": "Webhook processed successfully",
	})

}

// -----------------------------GET PAYMENT DETAILS------------------------------
func GetPaymentDetails(c *gin.Context) {

	orderID := c.Param("order_id")
	userID := c.MustGet("userID").(uint)

	var order models.Order

	if err := config.DB.
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {

		c.JSON(404, gin.H{
			"error": "Order not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"order_id":            order.ID,
		"payment_status":      order.PaymentStatus,
		"order_status":        order.Status,
		"total_amount":        order.TotalAmount,
		"razorpay_order_id":   order.RazorpayOrderID,
		"razorpay_payment_id": order.RazorpayPaymentID,
	})
}

func PlaceCODOrder(c *gin.Context) {

	orderID := c.Param("order_id")
	userID := c.MustGet("userID").(uint)

	var order models.Order

	if err := config.DB.
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {

		c.JSON(404, gin.H{
			"error": "Order not found",
		})
		return
	}

	if order.PaymentStatus == "paid" {
		c.JSON(400, gin.H{
			"error": "Order already paid",
		})
		return
	}

	order.PaymentMethod = "COD"
	order.PaymentStatus = "pending" // important for COD flow

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to update order",
		})
		return
	}

	c.JSON(200, gin.H{
		"order_id":       order.ID,
		"amount":         order.TotalAmount,
		"currency":       "INR",
		"payment_method": "COD",
		"payment_status": "pending",
	})

}
