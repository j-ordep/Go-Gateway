package domain

import "time"

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type Invoice struct {
	ID             string
	AccountId      string
	Status         string
	Description    string
	PaymentType    string
	CardLastDigits string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreditCard struct {
	
}