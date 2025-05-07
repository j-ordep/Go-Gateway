package service

import (
	"github.com/j-ordep/gateway/go-gateway/internal/domain"
)

// AccountService implementa a lógica de negócio para contas
// Aqui temos um exemplo de Inversão de Dependência, onde o service
// depende da interface do repository (AccountRepositoryInterface) em vez da implementação concreta.
// Isso permite:
// 1. Testar o service sem precisar de um banco de dados real
// 2. Trocar a implementação do repository sem afetar o service
// 3. Manter o service desacoplado da implementação do repository

type AccountService struct {
	repository domain.AccountRepositoryInterface
}

// NewAccountService é o construtor que recebe a dependência do repository
// Este é o ponto onde a injeção de dependência acontece no service
func NewAccountService(repository domain.AccountRepositoryInterface) *AccountService {
	return &AccountService{
		repository: repository,
	}
}
