package domain

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found") // é retornado quando a conta não é encontrada
	ErrDuplicatedAPIKey = errors.New("api key already exists") // é retornado quando há tentativa de criar conta com API key duplciada 
	ErrInvoiceNotFound = errors.New("invoice not found") // é retornado quando uma fatura não é encontrada
	ErrUnauthorizedAccess = errors.New("unauthorized access") // é retornado quando há tentativa de acesso não altorizado a um recurso
)