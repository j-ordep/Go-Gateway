package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	Name      string
	Email     string
	APIKey    string
	Balance   float64
	mu        sync.RWMutex // race conditions 
	CreatedAt time.Time
	UpdatedAt time.Time
}

func generateAPIKey() string {
	b := make([]byte, 16) // slice (parece ArrayList)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func NewAccount(name, email string) *Account {
	account := &Account {
		ID: uuid.New().String(),
		Name: name,
		Email: email,
		Balance: 0,
		APIKey: generateAPIKey(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return account
}

func (a *Account) AddBalance(amount float64) { // vulgo saldo
	a.mu.Lock() // espera a realização de uma transação para poder realizar outra
	defer a.mu.Unlock() // lembra do defer em JS, ele espera tudo executar antes de rodar
	a.Balance += amount
	a.UpdatedAt = time.Now()
}