package domain

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

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

	// gera um valor aleatorio, e verifica 
	randomSource := rand.New(rand.NewSource(time.Now().Unix()))

	var newStatus Status

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
	return nil
}