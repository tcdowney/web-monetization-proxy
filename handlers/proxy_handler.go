package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type ProxyHandler struct {
	BackendPort   int
	WalletPointer string
}

func (h *ProxyHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(fmt.Sprintf("http://0.0.0.0:%d", h.BackendPort))

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = BuildMonetizationResponseModifier(h.WalletPointer)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Host = url.Host

	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

	// proxy currently does not support gzipped/compressed responses from backends
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding
	req.Header.Set("Accept-Encoding", "identity")

	proxy.ServeHTTP(resp, req)
}

func BuildMonetizationResponseModifier(walletPointer string) func(*http.Response) error {
	return func(r *http.Response) error {
		if !contentTypeIsHTML(r) {
			log.Println("Content-Type is not HTML")
			return nil
		}

		doc, err := html.Parse(r.Body)
		if err != nil {
			log.Println(err)
			return nil
		}
		defer r.Body.Close()

		insertMonetizationMeta(doc, walletPointer)
		buf := bytes.NewBuffer([]byte{})
		html.Render(buf, doc)

		r.Body = ioutil.NopCloser(buf)
		r.Header["Content-Length"] = []string{fmt.Sprint(buf.Len())}
		return nil
	}
}

func insertMonetizationMeta(n *html.Node, walletPointer string) {
	if n.Type == html.ElementNode && n.Data == "head" {
		n.AppendChild(&html.Node{
			Type: html.ElementNode,
			Data: "meta",
			Attr: []html.Attribute{
				{Key: "name", Val: "monetization"},
				{Key: "content", Val: walletPointer},
			},
		})

		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		insertMonetizationMeta(c, walletPointer)
	}
}

func contentTypeIsHTML(r *http.Response) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "text/html")
}
