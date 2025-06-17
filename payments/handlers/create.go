package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"payments/pkg/json"
)

type createReq struct {
	FullName string  `json:"full_name" example:"John Doe"`
	Balance  float64 `json:"balance" example:"1000.00"`
} //@name CreateAccountRequest

// @Title Create Account
// @Description создает аккаунта
// @Tags Accounts Manage
// @Param   user_id  path  string  true  "ID пользователя"
// @Param   balance    body  createReq true  "Начальный счет аккаунта"
// @Success 200
// @Failure 400
// @Failure 500
// @Router  /account/create/{user_id} [post]
func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	req := &createReq{}
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

	//w.WriteHeader(http.StatusOK)
}
