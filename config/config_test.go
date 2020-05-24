package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/sclevine/spec"
	"github.com/tcdowney/web-monetization-proxy/config"
)

func TestConfig(t *testing.T) {
	spec.Run(t, "TestLoad", func(t *testing.T, when spec.G, it spec.S) {
		var cfg *config.Config

		it.Before(func() {
			err := os.Setenv("PROXY_PORT", "8081")
			if err != nil {
				t.Error(err)
			}
			err = os.Setenv("BACKEND_PORT", "9000")
			if err != nil {
				t.Error(err)
			}
			err = os.Setenv("WALLET_POINTER", "$wallet.example.com/🤑")
			if err != nil {
				t.Error(err)
			}

		})

		it("loads config from the environment", func() {
			var err error
			cfg, err = config.Load()
			if err != nil {
				t.Error(err)
			}

			if cfg.ProxyPort != 8081 {
				t.Errorf("Expected ProxyPort '%d' to match 8081", cfg.ProxyPort)
			}

			if cfg.BackendPort != 9000 {
				t.Errorf("Expected BackendPort '%d' to match 9000", cfg.BackendPort)
			}

			if cfg.WalletPointer != "$wallet.example.com/🤑" {
				t.Errorf("Expected WalletPointer '%s' to match '$wallet.example.com/🤑'", cfg.WalletPointer)
			}
		})

		when("PROXY_PORT is not provided", func() {
			it("defaults to 8080", func() {
				err := os.Unsetenv("PROXY_PORT")
				if err != nil {
					t.Error(err)
				}

				cfg, err = config.Load()
				if err != nil {
					t.Error(err)
				}

				if cfg.ProxyPort != 8080 {
					t.Errorf("Expected ProxyPort '%d' to match 8080", cfg.ProxyPort)
				}
			})
		})

		when("PROXY_PORT is not a parsable integer", func() {
			it("returns an error", func() {
				err := os.Setenv("PROXY_PORT", "ok")
				if err != nil {
					t.Error(err)
				}

				cfg, err = config.Load()
				if err == nil {
					t.Error("Expect error to have occurred")
				}

				if !strings.Contains(err.Error(), "integer") {
					t.Errorf("Expected error '%s' to say PROXY_PORT must be an integer", err)
				}
			})
		})

		when("BACKEND_PORT is not provided", func() {
			it("returns an error", func() {
				err := os.Unsetenv("BACKEND_PORT")
				if err != nil {
					t.Error(err)
				}

				cfg, err = config.Load()
				if err == nil {
					t.Error("Expect error to have occurred")
				}

				if !strings.Contains(err.Error(), "BACKEND_PORT is required") {
					t.Errorf("Expected error '%s' to say BACKEND_PORT is required", err)
				}
			})
		})

		when("BACKEND_PORT is not a parsable integer", func() {
			it("returns an error", func() {
				err := os.Setenv("BACKEND_PORT", "pls no")
				if err != nil {
					t.Error(err)
				}

				cfg, err = config.Load()
				if err == nil {
					t.Error("Expect error to have occurred")
				}

				if !strings.Contains(err.Error(), "integer") {
					t.Errorf("Expected error '%s' to say BACKEND_PORT must be an integer", err)
				}
			})
		})

		when("PROXY_PORT equals BACKEND_PORT", func() {
			it("returns an error", func() {
				err := os.Setenv("PROXY_PORT", "8081")
				if err != nil {
					t.Error(err)
				}
				err = os.Setenv("BACKEND_PORT", "8081")
				if err != nil {
					t.Error(err)
				}

				cfg, err = config.Load()
				if err == nil {
					t.Error("Expect error to have occurred")
				}

				if !strings.Contains(err.Error(), "PROXY_PORT cannot equal BACKEND_PORT") {
					t.Errorf("Expected error '%s' to say PROXY_PORT cannot equal BACKEND_PORT", err)
				}
			})
		})

		when("WALLET_POINTER is not provided", func() {
			it("returns an error", func() {
				err := os.Unsetenv("WALLET_POINTER")
				if err != nil {
					t.Error(err)
				}

				cfg, err = config.Load()
				if err == nil {
					t.Error("Expect error to have occurred")
				}

				if !strings.Contains(err.Error(), "WALLET_POINTER is required") {
					t.Errorf("Expected error '%s' to say WALLET_POINTER is required", err)
				}
			})
		})
	})
}
