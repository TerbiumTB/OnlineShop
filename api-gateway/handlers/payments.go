package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type PaymentsHandler struct {
	proxy *httputil.ReverseProxy
}

func NewPaymentsHandler() *PaymentsHandler {
	paymentsURL, _ := url.Parse(os.Getenv("PAYMENTS_URL"))
	paymentsProxy := httputil.NewSingleHostReverseProxy(paymentsURL)
	return &PaymentsHandler{paymentsProxy}
}

func (h *PaymentsHandler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(wr, r)
}
