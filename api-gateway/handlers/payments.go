package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var paymentsProxy *httputil.ReverseProxy

func init() {
	paymentsURL, _ := url.Parse(os.Getenv("PAYMENTS_URL"))
	paymentsProxy = httputil.NewSingleHostReverseProxy(paymentsURL)
}

func PaymentsHandler(wr http.ResponseWriter, r *http.Request) {
	paymentsProxy.ServeHTTP(wr, r)
}
