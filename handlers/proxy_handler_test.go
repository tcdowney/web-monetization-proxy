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
		var proxyHandler handlers.ProxyHandler
		var addWebMonetizationMetaFunc func(*http.Response) error

		it.Before(func() {
			proxyHandler = handlers.ProxyHandler{
				BackendPort:   1337,
				WalletPointer: "$wallet.example.com/🤑",
			}
			addWebMonetizationMetaFunc = handlers.BuildMonetizationResponseModifier(proxyHandler.WalletPointer)
		})

		when("the response is not HTML", func() {
			it.Before(func() {
				response = &http.Response{
					StatusCode: 200,
					ProtoMajor: 1,
					ProtoMinor: 1,
					Header:     make(http.Header),
					Body:       ioutil.NopCloser(bytes.NewBufferString("console.log('hello world')")),
				}

				response.Header.Set("Content-Type", mime.TypeByExtension(".js"))
			})

			it("does not modify the response", func() {
				if err := addWebMonetizationMetaFunc(response); err != nil {
					t.Error(err)
				}

				bodyBytes, err := ioutil.ReadAll(response.Body)
				if err != nil {
					t.Error(err)
				}

				bodyString := string(bodyBytes)
				if strings.Contains(bodyString, "<meta name=\"monetization\" content=\"$wallet.example.com/🤑\"/>") {
					t.Error(fmt.Sprintf("<meta> tag was added: %s", bodyString))
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
					if err := addWebMonetizationMetaFunc(response); err != nil {
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
					if err := addWebMonetizationMetaFunc(response); err != nil {
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
			})
		})
	})
}
