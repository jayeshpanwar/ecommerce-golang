package config

import (
	"os"

	"github.com/razorpay/razorpay-go"
)

func NewRazorpayClient() *razorpay.Client {
	return razorpay.NewClient(
		os.Getenv("RAZORPAY_KEY_ID"),
		os.Getenv("RAZORPAY_KEY_SECRET"),
	)
}
