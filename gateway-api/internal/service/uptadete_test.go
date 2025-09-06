package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/j-ordep/gateway/go-gateway/internal/domain"
	"github.com/j-ordep/gateway/go-gateway/internal/service"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Save(account *domain.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) FindByAPIKey(apiKey string) (*domain.Account, error) {
	args := m.Called(apiKey)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) FindById(id string) (*domain.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) UpdateBalance(account *domain.Account) error {
	args := m.Called(account)
	return args.Error(1)
}

func TestUpdateBalance(t *testing.T) {
	mockRepo := &MockAccountRepository{}
	service := service.NewAccountService(mockRepo)
    
    // GIVEN - Preparar dados
	existingAccount := domain.NewAccount("João", "joao@gmail.com")

	// - cria uma variavel que pega a apikey de account acima ( variavel := existingAccount.APIKey)
	apiKey := existingAccount.APIKey

	// - cria uma variavel que pega o balance de account acima ( variavel := existingAccount.Balance)
	initialBalance := existingAccount.Balance

    amountToAdd := 100.00

    // WHEN - Configurar comportamento do mock (.On .Return)
	mockRepo.On("FindByAPIKey", apiKey).Return(existingAccount, nil)
	mockRepo.On("UpdateBalance", existingAccount).Return(nil)
    
    // 4. ACT - Executar o método
	result, err := service.UpdateBalance(apiKey, amountToAdd)
    
    // 5. ASSERT - Verificar resultados
    assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, initialBalance + amountToAdd, result.Balance)
    // 6. VERIFY - Verificar se métodos foram chamados
	mockRepo.AssertExpectations(t)
}