package domain

type AccountRepository interface {
	Save(account *Account) error
	FindByAPIKey(apiKey string) (*Account, error)
	UpdateBalance(account *Account) error
}

type InvoiceRepository interface {
	Save(invoice *Invoice) error
	FindById(id string) (*Invoice, error)
	FindByAccountId(id string) ([]*Invoice, error)
	UpdateStatus(status Status) error
}