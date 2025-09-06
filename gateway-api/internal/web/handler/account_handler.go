package handler

import (
	"encoding/json"
	"net/http"

	"github.com/j-ordep/gateway/go-gateway/internal/dto"
	"github.com/j-ordep/gateway/go-gateway/internal/service"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateAccountInput
	
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	output, err := h.accountService.CreateAccount(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(output)

}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-KEY")
	if apiKey == "" {
		http.Error(w, "X-API-KEY header is required", http.StatusUnauthorized)
		return 
	}

	output, err := h.accountService.FindByAPIKey(apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(output)
	
}