package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tdb/ms-payment/src/cmd/enum"
	"github.com/tdb/ms-payment/src/cmd/services"
)

type PaymentResponse struct {
	Message string             `json:"message"`
	Status  enum.PaymentStatus `json:"status"`
}

type PaymentHandler interface {
	RegisterPayment(w http.ResponseWriter, r *http.Request)
}

type paymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) PaymentHandler {
	return &paymentHandler{
		paymentService: paymentService,
	}
}

func (ph *paymentHandler) RegisterPayment(w http.ResponseWriter, r *http.Request) {
	var paymentInfo services.PaymentInput

	err := json.NewDecoder(r.Body).Decode(&paymentInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ph.paymentService.RegisterPayment(paymentInfo)
	if err != nil {
		fmt.Print(err)
		if err.Error() == "account not found" {
			response := PaymentResponse{
				Message: "Account not found",
				Status:  enum.DENIED,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if err.Error() == "insufficient balance" {
			response := PaymentResponse{
				Message: "Insufficient balance",
				Status:  enum.DENIED,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		if err.Error() == "invalid payment method" {
			response := PaymentResponse{
				Message: "Invalid payment method",
				Status:  enum.DENIED,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		response := PaymentResponse{
			Message: "Error registering payment",
			Status:  enum.DENIED,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := PaymentResponse{
		Message: "Payment registered successfully",
		Status:  enum.APPROVED,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
