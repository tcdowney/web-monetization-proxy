package handlers_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"testing"

	"github.com/sclevine/spec"
	"github.com/tcdowney/web-monetization-proxy/handlers"
)

func TestAddWebMonetizationMeta(t *testing.T) {
	spec.Run(t, "TestBuildMonetizationResponseModifier", func(t *testing.T, when spec.G, it spec.S) {
		var response *http.Response
		var proxyHandler *handlers.ProxyHandler
		var proxyScriptHandler *handlers.ProxyHandler
		var addWmMetaFunc func(*http.Response) error
		var addWmMetaAndScriptFunc func(*http.Response) error
		expectedScript :=
			`<script>document.monetization&&document.monetization.addEventListener("monetizationprogress",e=>{const{receipt:t}=e.detail;null!==t&&fetch("https://verifier.com/balances/123:creditReceipt",{method:"POST",body:t})});</script>`

		it.Before(func() {
			proxyHandler = handlers.NewProxyHandler(1337, "$wallet.example.com/🤑", "")
			addWmMetaFunc = handlers.BuildMonetizationResponseModifier(proxyHandler.PaymentPointer, proxyHandler.ReceiptSubmitter)
			proxyScriptHandler = handlers.NewProxyHandler(1337, "$wallet.example.com/🤑", "https://verifier.com/balances/123:creditReceipt")
			addWmMetaAndScriptFunc = handlers.BuildMonetizationResponseModifier(proxyScriptHandler.PaymentPointer, proxyScriptHandler.ReceiptSubmitter)
		})

		when("the response is not HTML", func() {
			bodyString := "console.log('hello world')"
			it.Before(func() {
				response = &http.Response{
					StatusCode: 200,
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     make(http.Header),
					Body:       ioutil.NopCloser(bytes.NewBufferString(bodyString)),
				}

				response.Header.Set("Content-Type", mime.TypeByExtension(".js"))
			})

			it("does not modify the response", func() {
				if err := addWmMetaFunc(response); err != nil {
					t.Error(err)
				}

				bodyBytes, err := ioutil.ReadAll(response.Body)
				if err != nil {
					t.Error(err)
				}

				if string(bodyBytes) != bodyString {
					t.Error(fmt.Sprintf("response was modified: %s", bodyString))
				}
			})

			it("does not modify the response with a script", func() {
				if err := addWmMetaAndScriptFunc(response); err != nil {
					t.Error(err)
				}

				bodyBytes, err := ioutil.ReadAll(response.Body)
				if err != nil {
					t.Error(err)
				}

				if string(bodyBytes) != bodyString {
					t.Error(fmt.Sprintf("response was modified: %s", bodyString))
				}
			})
		})

		when("the response is HTML", func() {
			when("the response HTML has a <head> tag", func() {
				it.Before(func() {
					response = &http.Response{
						StatusCode: 200,
						ProtoMajor: 1,
						ProtoMinor: 1,
						Header:     make(http.Header),
						Body:       ioutil.NopCloser(bytes.NewBufferString("<html><head></head></html>")),
					}

					response.Header.Set("Content-Type", mime.TypeByExtension(".html"))
				})

				it("adds the monetization <meta> tag", func() {
					if err := addWmMetaFunc(response); err != nil {
						t.Error(err)
					}

					bodyBytes, err := ioutil.ReadAll(response.Body)
					if err != nil {
						t.Error(err)
					}

					bodyString := string(bodyBytes)
					if !strings.Contains(bodyString, "<head><meta name=\"monetization\" content=\"$wallet.example.com/🤑\"/></head>") {
						t.Error(fmt.Sprintf("<meta> tag not added: %s", bodyString))
					}
				})

				it("adds the monetization <meta> tag and receipt submission <script>", func() {
					if err := addWmMetaAndScriptFunc(response); err != nil {
						t.Error(err)
					}

					bodyBytes, err := ioutil.ReadAll(response.Body)
					if err != nil {
						t.Error(err)
					}

					bodyString := string(bodyBytes)
					if !strings.Contains(bodyString, "<meta name=\"monetization\" content=\"$wallet.example.com/🤑\"/>") {
						t.Error(fmt.Sprintf("<meta> tag not added: %s", bodyString))
					}

					if !strings.Contains(bodyString, expectedScript) {
						t.Error(fmt.Sprintf("<script> tag not added: %s", bodyString))
					}
				})
			})

			when("the response HTML does NOT have a <head> tag", func() {
				it.Before(func() {
					response = &http.Response{
						StatusCode: 200,
						ProtoMajor: 1,
						ProtoMinor: 1,
						Header:     make(http.Header),
						Body:       ioutil.NopCloser(bytes.NewBufferString("<html><body><p>ok</p></body</html>")),
					}

					response.Header.Set("Content-Type", mime.TypeByExtension(".html"))
				})

				it("adds the monetization <meta> tag and <head> tag", func() {
					if err := addWmMetaFunc(response); err != nil {
						t.Error(err)
					}

					bodyBytes, err := ioutil.ReadAll(response.Body)
					if err != nil {
						t.Error(err)
					}

					bodyString := string(bodyBytes)
					if !strings.Contains(bodyString, "<head><meta name=\"monetization\" content=\"$wallet.example.com/🤑\"/></head>") {
						t.Error(fmt.Sprintf("<meta> tag not added: %s", bodyString))
					}
				})

				it("adds the monetization <meta>, receipt submission <script>, and <head> tags", func() {
					if err := addWmMetaAndScriptFunc(response); err != nil {
						t.Error(err)
					}

					bodyBytes, err := ioutil.ReadAll(response.Body)
					if err != nil {
						t.Error(err)
					}

					bodyString := string(bodyBytes)
					if !strings.Contains(bodyString, "<meta name=\"monetization\" content=\"$wallet.example.com/🤑\"/>") {
						t.Error(fmt.Sprintf("<meta> tag not added: %s", bodyString))
					}

					if !strings.Contains(bodyString, expectedScript) {
						t.Error(fmt.Sprintf("<script> tag not added: %s", bodyString))
					}
				})
			})
		})
	})
}
