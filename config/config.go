package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ProxyPort            int
	BackendPort          int
	PaymentPointer       string
	ReceiptSubmissionUrl string
}

func Load() (*Config, error) {
	proxyPortString := os.Getenv("PROXY_PORT")
	if proxyPortString == "" {
		proxyPortString = "8080"
	}

	proxyPort, err := strconv.Atoi(proxyPortString)
	if err != nil {
		return nil, fmt.Errorf("PROXY_PORT '%s' must be an integer", proxyPortString)
	}

	backendPortString := os.Getenv("BACKEND_PORT")
	if backendPortString == "" {
		return nil, fmt.Errorf("BACKEND_PORT is required")
	}

	backendPort, err := strconv.Atoi(backendPortString)
	if err != nil {
		return nil, fmt.Errorf("BACKEND_PORT '%s' must be an integer", backendPortString)
	}

	if proxyPort == backendPort {
		return nil, fmt.Errorf("PROXY_PORT cannot equal BACKEND_PORT")
	}

	paymentPointer := os.Getenv("PAYMENT_POINTER")
	if paymentPointer == "" {
		return nil, fmt.Errorf("PAYMENT_POINTER is required")
	}

	receiptSubmissionUrl := os.Getenv("RECEIPT_SUBMISSION_URL")

	c := &Config{
		ProxyPort:            proxyPort,
		BackendPort:          backendPort,
		PaymentPointer:       paymentPointer,
		ReceiptSubmissionUrl: receiptSubmissionUrl,
	}

	return c, nil
}
