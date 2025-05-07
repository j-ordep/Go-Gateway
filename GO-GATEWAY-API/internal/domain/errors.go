package domain

import "errors"

var (
	// é retornado quando a conta não é encontrada
	ErrAccountNotFound = errors.New("account not found") 

	// é retornado quando há tentativa de criar conta com API key duplciada 
	ErrDuplicatedAPIKey = errors.New("api key already exists") 

	// é retornado quando uma fatura não é encontrada
	ErrInvoiceNotFound = errors.New("invoice not found") 

	// é retornado quando há tentativa de acesso não altorizado a um recurso
	ErrUnauthorizedAccess = errors.New("unauthorized access") 
)