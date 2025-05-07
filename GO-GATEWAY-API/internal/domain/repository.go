package domain

type AccountRepositoryInterface interface {
	Save(account *Account) error
	FindByAPIKey(apiKey string) (*Account, error)
	FindById(id string) (*Account, error)
	Update(account *Account) error
} 
