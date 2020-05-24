package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ProxyPort     int
	BackendPort   int
	WalletPointer string
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

	walletPointer := os.Getenv("WALLET_POINTER")
	if walletPointer == "" {
		return nil, fmt.Errorf("WALLET_POINTER is required")
	}

	c := &Config{
		ProxyPort:     proxyPort,
		BackendPort:   backendPort,
		WalletPointer: walletPointer,
	}

	return c, nil
}
