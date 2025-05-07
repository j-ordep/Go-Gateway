package dto

import (
	"time"

	"github.com/j-ordep/gateway/go-gateway/internal/domain"
)

type CreateAccountInput struct {
	Name  string `json:"name"`
	Email string `jason:"email"`
}

type AccountOuput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	APIKey    string    `json:"api_key,omitempty"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToAccount(input CreateAccountInput) *domain.Account { // dto não tem ponteiros, ele transfere os dados e depois jogado fora pelo GC
	return domain.NewAccount(input.Name, input.Email)
}

func FromAccount(account *domain.Account) AccountOuput {
	return AccountOuput{
		ID: account.ID,
		Name: account.Name,
		Email: account.Email,
		Balance: account.Balance,
		APIKey: account.APIKey,
		CreatedAt: account.CreatedAt ,
		UpdatedAt: account.UpdatedAt,
	}
}