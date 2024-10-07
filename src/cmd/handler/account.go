package handler

import (
	"encoding/json"
	"net/http"

	"github.com/tdb/ms-payment/src/cmd/services"
)

type AccountHandler interface {
	CreateAccount(w http.ResponseWriter, r *http.Request)
}

type accountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(accountService services.AccountService) AccountHandler {
	return &accountHandler{
		accountService: accountService,
	}
}

func (h *accountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var accountInfo services.AccountInput

	err := json.NewDecoder(r.Body).Decode(&accountInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.accountService.CreateAccount(accountInfo)
	if err != nil {
		if err.Error() == "pix key or credit card number is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err.Error() == "name is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
