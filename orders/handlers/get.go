package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"orders/pkg/json"
)

// @Title Get All Orders
// @Description Получить все заказы
// @Tags Order Info
// @Produce json
// @Success 200 {array} model.Order "Информация о заказе"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /get [get]
func (h *Handler) AllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.s.All()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.ToJSON(orders, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Title Get Order By ID
// @Description Получить заказ по ID
// @Tags Order Info
// @Produce json
// @Param   id  path  string  true  "ID заказа"
// @Success 200 {object} model.Order "Информация о заказе"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /get/{id} [get]
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	order, err := h.s.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.ToJSON(order, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
