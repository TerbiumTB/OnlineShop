package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"payments/pkg/json"
)

type createAccountRequest struct {
	FullName string  `json:"full_name"`
	Balance  float64 `json:"balance"`
	//Descr  string  `json:"descr"`
}

// @Title Create Account
// @Description создает аккаунта
// @Tags Аккаунты
// @Param   user_id  path  string  true  "ID пользователя"
// @Param   balance    body  createAccountRequest true  "Начальный счет аккаунта"
// @Success 200
// @Failure 400
// @Failure 500
// @Router  /account/create/{user_id} [post]
func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	req := &createAccountRequest{}
	err := json.FromJSON(req, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.as.Add(userID, req.FullName, req.Balance)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
