package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// var ordersProxy *httputil.ReverseProxy
type OrdersHandler struct {
	proxy *httputil.ReverseProxy
}

func NewOrdersHandler() *OrdersHandler {
	ordersURL, _ := url.Parse(os.Getenv("ORDERS_URL"))
	ordersProxy := httputil.NewSingleHostReverseProxy(ordersURL)
	return &OrdersHandler{ordersProxy}
}

func (h *OrdersHandler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(wr, r)
}
