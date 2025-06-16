package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"orders/pkg/json"
)

type request struct {
	Amount float64 `json:"amount"`
	Descr  string  `json:"descr"`
}

// @Title CreateOrder
// @Description Оформляет заказ
// @Tags Загрузка
// @Param   user_id  path  string  true  "ID пользователя"
// @Param   order     body  request true  "Детали заказа"
// @Success 200
// @Failure 400
// @Failure 500
// @Router  /order/create/{user_id} [post]
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	h.l.Printf("Handling create order request for user %s", userID)
	req := &request{}
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

	w.WriteHeader(http.StatusOK)
}
