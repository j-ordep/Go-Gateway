package repository

import (
	"database/sql"
	"time"

	"github.com/j-ordep/gateway/go-gateway/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (repo *AccountRepository) Save(account *domain.Account) error {

    // Prepare: Cria uma query preparada que pode ser executada múltiplas vezes
    // - Previne SQL Injection (parametrização automática)
    // - Melhora performance em execuções repetidas (query é compilada uma vez)
    // - Validação da sintaxe SQL antes da execução
    // Neste caso, usamos Prepare pois não precisamos de transação ou lock de linha
	stmt, err := repo.db.Prepare(`
		INSERT INTO accounts (id, name, email, api_key, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)	
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		account.ID,
		account.Name,
		account.Email,
		account.APIKey,
		account.Balance,
		account.CreatedAt,
		account.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *AccountRepository) FindByAPIKey(apiKey string) (*domain.Account, error) {
	query := `
		SELECT id, name, email, api_key, balance, created_at, updated_at
		FROM accounts
		WHERE api_key = $1
	`

	// Poderíamos fazer o Scan diretamente em &account.CreatedAt e &account.UpdatedAt,
	// porém, por segurança e para evitar problemas de tipo caso a struct Account mude (ex: ponteiros, tipos customizados),
	// utilizamos variáveis intermediárias.
	var createdAt, updatedAt time.Time
	var account domain.Account

	err := repo.db.QueryRow(query, apiKey).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.APIKey,
		&account.Balance,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	account.CreatedAt = createdAt
	account.UpdatedAt = updatedAt

	return &account, nil
}

func (repo *AccountRepository) UpdateBalance(account *domain.Account) error {
	
	// 1. FOR UPDATE: Permite fazer lock pessimista na linha (impede leituras/escritas concorrentes)
    // 2. Atomicidade: Garante que SELECT + UPDATE aconteçam como uma operação única
    // 3. Isolamento: Previne race condition em atualizações simultâneas de saldo
    // 4. Consistência: Se algo falhar, Rollback() desfaz todas as alterações
	tx, err := repo.db.Begin() 
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	var currentBalance float64

	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, account.ID).Scan(&currentBalance)
	if err == sql.ErrNoRows {
		return domain.ErrAccountNotFound
	}
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE accounts
		SET balance = $1, updated_at = $2
		WHERE id = $3
	`, account.Balance, time.Now(), account.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}