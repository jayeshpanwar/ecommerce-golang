package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendOrderConfirmationEmail(
	toEmail string,
	orderID uint,
	total float64,
) error {

	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth(
		"",
		from,
		password,
		host,
	)

	subject := fmt.Sprintf(
		"Order #%d Confirmed",
		orderID,
	)

	body := fmt.Sprintf(
		`Thank you for shopping with us!

Order ID: %d
Total Amount: ₹%.2f

Your order has been placed successfully.

We will notify you once the seller accepts and ships your order.`,
		orderID,
		total,
	)

	message := []byte(
		"Subject: " + subject + "\r\n" +
			"\r\n" +
			body,
	)

	return smtp.SendMail(
		host+":"+port,
		auth,
		from,
		[]string{toEmail},
		message,
	)
}
