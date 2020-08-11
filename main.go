package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tcdowney/web-monetization-proxy/config"
	"github.com/tcdowney/web-monetization-proxy/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	proxyHandler := handlers.NewProxyHandler(cfg.BackendPort, cfg.PaymentPointer, cfg.ReceiptSubmissionUrl)

	mux := http.NewServeMux()
	mux.Handle("/", proxyHandler)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.ProxyPort), mux)
}
