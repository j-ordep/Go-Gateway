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
	query := `
		SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at
		FROM invoices
		WHERE id = $1
	`

	var invoice domain.Invoice

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

// pode retornar varios Invoices, pois varios invoices podem ter o mesmo accountID, (1 account pode ter mais de um invoice)
func (r *InvoiceRepository) FindByAccountId(accountId string) ([]*domain.Invoice, error) {
	query := `
		SELECT id, account_id, amount, status, description, payment_type, card_last_digits, created_at, updated_at
		FROM invoices
		WHERE account_id = $1
	`
	
	rows, err := r.db.Query(query, accountId)
	if err != nil {
		return nil, err
	}

	var invoices []*domain.Invoice
	
	for rows.Next() {
		var invoice domain.Invoice
		err := rows.Scan(
			&invoice.ID,
			&invoice.AccountId,
			&invoice.Amount,
			&invoice.Status,
			&invoice.Description,
			&invoice.PaymentType,
			&invoice.CardLastDigits,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		invoices = append(invoices, &invoice)
	}
	defer rows.Close()

	return invoices, nil
}

// unica reponsabilidade do repository é salva no DB, o invoice já vem alterado
func (r *InvoiceRepository) UpdateStatus(invoice *domain.Invoice) error {
	query := `
		UPDATE invoices 
		SET status = $1, updated_at = $2 
		WHERE id = $3
	`

	rows, err := r.db.Exec(query, invoice.Status, invoice.UpdatedAt, invoice.ID)
	if err != nil {
		return err
	}

	rowsAffedted, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffedted == 0 {
		return domain.ErrInvoiceNotFound
	}

	return nil
}