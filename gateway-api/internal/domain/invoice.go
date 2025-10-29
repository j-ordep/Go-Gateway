package domain

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// por que usar Status como um tipo?
// - type safety (Segurança de tipo)
// - Centralização dos valores válidos
// - Facilita validações e refatorações
type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type Invoice struct {
	ID             string
	AccountId      string
	Amount         float64
	Status         Status
	Description    string
	PaymentType    string
	CardLastDigits string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreditCard struct {
	Number string
	CVV string
	ExpiryMonth int
	ExpiryYear int
	CardholderName string
}

func NewInvoice(accountId string, amount float64, description string, paymentType string, card CreditCard) (*Invoice, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// len(card.Number) = 16
	lastDigits := card.Number[len(card.Number)-4:] // 16 - 4 = [12:] (basicamente ele pega do 12º numero para frente, ou seja ultimos 4 numeros)
	
	return &Invoice{
		ID: uuid.New().String(),
		AccountId: accountId,
		Amount: amount,
		Status: StatusPending,
		Description: description,
		PaymentType: paymentType,
		CardLastDigits: lastDigits,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (i *Invoice) Process() error {
	if i.Amount > 10000 {
		return nil // mantem o status como pendente (StatusPending), com isso enviamos o invoice para apache kafka
	}

	// gera um valor numerico aleatorio de 0.0 até 1.0 (ex: 90.12, 0.45, 0.83)
	// Isso garante que a cada execução o gerador produza sequências diferentes de números aleatórios.
	// Se usássemos sempre a mesma semente, os resultados seriam sempre iguais, o que não é desejado para simular aleatoriedade real.
	randomSource := rand.New(rand.NewSource(time.Now().Unix()))

	var newStatus Status

	// 70% chance de aprovação
	if randomSource.Float64() <= 0.7 {
		newStatus = StatusApproved
	} else {
		newStatus = StatusRejected
	}

	i.Status = newStatus

	return nil
}

func (i *Invoice) UpdateStatus(newStatus Status) error {
	if i.Status != StatusPending {
		return ErrInvalidStatus
	}
	
	i.Status = newStatus
	i.UpdatedAt = time.Now()
	return nil
}