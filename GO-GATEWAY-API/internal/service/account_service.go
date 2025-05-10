package service

import (
	"github.com/j-ordep/gateway/go-gateway/internal/domain"
	"github.com/j-ordep/gateway/go-gateway/internal/dto"
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

func NewAccountService(repository domain.AccountRepositoryInterface) *AccountService {
	return &AccountService{repository: repository,}
}

func (s *AccountService) CreateAccount(input dto.CreateAccountInput) (*dto.AccountOuput, error) {
	account := dto.ToAccount(input)

	existingAccount, err := s.repository.FindByAPIKey(account.APIKey)
	if err != nil && err != domain.ErrAccountNotFound {
		return nil, err
	}
	if existingAccount != nil {
		return nil, domain.ErrDuplicatedAPIKey
	}
	err = s.repository.Save(account)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output, nil
}

func (s *AccountService) UpdateBalance(apiKey string, amount float64) (*dto.AccountOuput, error) {

	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}
	
	account.AddBalance(amount)
	err = s.repository.UpdateBalance(account)
	if err != nil {
		return nil, err
	}

	output := dto.FromAccount(account)
	return &output,nil
}

func (s *AccountService) FindByAPIKey(apiKey string) (*dto.AccountOuput, error) {
	account, err := s.repository.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}
	output := dto.FromAccount(account)
	return &output, nil
}

func (s *AccountService) FindById(id string) (*dto.AccountOuput, error) {
	account, err := s.repository.FindById(id)
	if err != nil {
		return nil, err
	}
	output := dto.FromAccount(account)
	return &output, nil
}
