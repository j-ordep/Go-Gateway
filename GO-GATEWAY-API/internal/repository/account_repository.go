package repository

import (
	"database/sql"
	"time"

	"github.com/j-ordep/gateway/go-gateway/internal/domain"
)

// AccountRepository é a implementação concreta da interface AccountRepositoryInterface
// Aqui temos um exemplo de Injeção de Dependência, onde o AccountRepository
// recebe sua dependência (db *sql.DB) através do construtor, em vez de criar internamente.
// Isso permite:
// 1. Testabilidade (podemos injetar um mock do banco)
// 2. Flexibilidade (podemos trocar o banco facilmente)
// 3. Desacoplamento (o repositório não precisa saber como criar a conexão)

type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository é o construtor que recebe a dependência
// Este é o ponto onde a injeção de dependência acontece
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (repo *AccountRepository) Save(account *domain.Account) error {

	// statment
	stmt, err := repo.db.Prepare(`
		INSERT INTO accounts (id, name, email, api_key, balance, create_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)	
	`)
	if err != nil {
		return err
	}
	defer stmt.Close() // defer = eperar tudo rodar ai você executa o Close

	// substitui adicionando os elementos no VALUES da query
	_, err = stmt.Exec(
		account.ID,
		account.Name,
		account.Email,
		account.APIKey,
		account.CreatedAt,
		account.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *AccountRepository) FindByAPIKey(apiKey string) (*domain.Account, error) {
	var account domain.Account
	var createdAt, updatedAt time.Time

	err := repo.db.QueryRow(`
		SELECT id, name, email, api_key, balance, created_at, update_at
		FROM accounts
		WHERE api_key = $1
	`, apiKey).Scan( // $1 é onde vai a apiKey (que entrou como parametro)
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

	tx, err := repo.db.Begin() // tx = transaction
	if err != nil {
		return err
	}

	defer tx.Rollback() // Se algo der errado, desfaz todas as operações

	var currentBalance float64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`,
		account.ID).Scan(&currentBalance)

	if err == sql.ErrNoRows {
		return domain.ErrAccountNotFound
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		UPDATE accounts
		SET balance = $1, update_at = $2
		WHERE id = $3
	`, account.Balance, time.Now(), account.ID)
	if err != nil {
		return err
	}
	return tx.Commit()

}