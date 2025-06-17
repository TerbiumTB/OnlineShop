package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"orders/pkg/json"
)

type createReq struct {
	Amount float64 `json:"amount" example:"500.00"`
	Descr  string  `json:"descr" example:"Labuba"`
} //@name CreateOrderRequest

// @Title Create Order
// @Description Оформляет заказ
// @Tags Order Manage
// @Param   user_id  path  string  true  "ID пользователя"
// @Param   order     body  createReq true  "Детали заказа"
// @Success 200
// @Failure 400
// @Failure 500
// @Router  /create/{user_id} [post]
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]

	req := &createReq{}
	err := json.FromJSON(req, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.Add(userID, req.Amount, req.Descr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
