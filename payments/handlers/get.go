package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"payments/pkg/json"
)

// @Title Get all accounts
// @Description Возвращает все счета
// @Tags Accounts Info
// @Produce json
// @Success 200 {array} model.Account
// @Router  /account/get [get]
func (h *Handler) AllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.as.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.ToJSON(accounts, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Title Get Account By ID
// @Description Получить счет по ID
// @Tags Accounts Info
// @Produce json
// @Param   id  path  string  true  "ID пользователя"
// @Success 200 {object} model.Account "Информация о счете"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /account/get/{user_id} [get]
func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["user_id"]
	account, err := h.as.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.ToJSON(account, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//w.WriteHeader(http.StatusOK)
}
