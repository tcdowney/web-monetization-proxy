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
	BackendPort      int
	PaymentPointer   string
	ReceiptSubmitter string
}

func NewProxyHandler(backendPort int, paymentPointer string, receiptSubmissionUrl string) *ProxyHandler {
	var receiptSubmitter string
	if receiptSubmissionUrl != "" {
		receiptSubmitter = fmt.Sprintf(`document.monetization&&document.monetization.addEventListener("monetizationprogress",e=>{const{receipt:t}=e.detail;null!==t&&fetch("%s",{method:"POST",body:t})});`, receiptSubmissionUrl)
	}
	return &ProxyHandler{
		BackendPort:      backendPort,
		PaymentPointer:   paymentPointer,
		ReceiptSubmitter: receiptSubmitter,
	}
}

func (h *ProxyHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(fmt.Sprintf("http://0.0.0.0:%d", h.BackendPort))

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ModifyResponse = BuildMonetizationResponseModifier(h.PaymentPointer, h.ReceiptSubmitter)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Host = url.Host

	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

	// proxy currently does not support gzipped/compressed responses from backends
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding
	req.Header.Set("Accept-Encoding", "identity")

	proxy.ServeHTTP(resp, req)
}

func BuildMonetizationResponseModifier(paymentPointer string, receiptSubmitter string) func(*http.Response) error {
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

		insertInHead(doc, paymentPointer, receiptSubmitter)
		buf := bytes.NewBuffer([]byte{})
		html.Render(buf, doc)

		r.Body = ioutil.NopCloser(buf)
		r.Header["Content-Length"] = []string{fmt.Sprint(buf.Len())}
		return nil
	}
}

func insertInHead(n *html.Node, paymentPointer string, receiptSubmitter string) {
	if n.Type == html.ElementNode && n.Data == "head" {
		insertMonetizationMeta(n, paymentPointer)
		if receiptSubmitter != "" {
			insertReceiptSubmitter(n, receiptSubmitter)
		}

		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		insertInHead(c, paymentPointer, receiptSubmitter)
	}
}

func insertMonetizationMeta(n *html.Node, paymentPointer string) {
	if n.Type == html.ElementNode && n.Data == "head" {
		n.AppendChild(&html.Node{
			Type: html.ElementNode,
			Data: "meta",
			Attr: []html.Attribute{
				{Key: "name", Val: "monetization"},
				{Key: "content", Val: paymentPointer},
			},
		})
	}
}

func insertReceiptSubmitter(n *html.Node, receiptSubmitter string) {
	if n.Type == html.ElementNode && n.Data == "head" {
		script := html.Node{
			Type: html.ElementNode,
			Data: "script",
		}
		script.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: receiptSubmitter,
		})
		n.AppendChild(&script)
	}
}

func contentTypeIsHTML(r *http.Response) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "text/html")
}
