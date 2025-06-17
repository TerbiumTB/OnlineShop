package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"payments/pkg/json"
)

type updateReq struct {
	Amount float64 `json:"amount" example:"100.00"`
} //@name UpdateAccountRequest

// @Title Update Account
// @Description изменить баланс счета
// @Tags Accounts Manage
// @Param   user_id  path  string  true  "ID пользователя"
// @Param   update_request   body  updateReq true  "Изменение счета"
// @Success 200
// @Failure 400
// @Failure 500
// @Router  /account/update/{user_id} [patch]
func (h *Handler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	req := &updateReq{}
	err := json.FromJSON(req, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.as.Update(userID, req.Amount)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
}
