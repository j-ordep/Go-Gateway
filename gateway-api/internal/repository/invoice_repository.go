package repository

import (
	"database/sql"

	"github.com/j-ordep/gateway/go-gateway/internal/domain"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db:db}
}

func (r *InvoiceRepository) Save(invoice *domain.Invoice) error {
	query := `
		INSERT INTO invoices (id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, &7, $8, $9)
	`

	_, err := r.db.Exec(query, 
		invoice.ID, 
		invoice.AccountId, 
		invoice.Amount, 
		invoice.Status, 
		invoice.Description, 
		invoice.PaymentType, 
		invoice.CardLastDigits, 
		invoice.CreatedAt, 
		invoice.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *InvoiceRepository) FindById(id string) (*domain.Invoice, error) {
	var invoice domain.Invoice

	query := `
		SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at
		FROM invoices
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&invoice.AccountId, 
		&invoice.Amount, 
		&invoice.Status, 
		&invoice.Description, 
		&invoice.PaymentType, 
		&invoice.CardLastDigits, 
		&invoice.CreatedAt, 
		&invoice.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrInvoiceNotFound
	}

	if err != nil {
		return nil, err
	}

	return &invoice ,nil
}