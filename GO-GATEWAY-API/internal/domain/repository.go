package domain

type AccountRepositoryInterface interface {
	Save(account *Account)
	FindByAPIKey(apiKey string) (*Account, error)
	FindById(id string) (*Account, error)
	Update(account *Account) error
} 
