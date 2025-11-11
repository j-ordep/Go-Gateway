package service

import (
	"github.com/j-ordep/gateway/go-gateway/internal/domain"
	"github.com/j-ordep/gateway/go-gateway/internal/dto"
)

type InvoiceService struct {
	invoiceRepository domain.InvoiceRepository
	accountService    AccountService
}

func NewInvoiceService(invoiceRepository domain.InvoiceRepository, accountService AccountService) *InvoiceService {
	return &InvoiceService{
		invoiceRepository: invoiceRepository,
		accountService:    accountService,
	}
}

func (s *InvoiceService) Create(input dto.CreateInvoiceInput) (*dto.InvoiceOutput, error) {
	
	// 1. achar o dono da fatura (accountId)
	accountOutput, err := s.accountService.FindByAPIKey(input.APIKey)
	if err != nil {
		return nil, err
	}

	// 2. tranformar o input (dto - invoice + credit_card) em invoice domain
	invoice, err := dto.ToInvoice(input, accountOutput.ID)
	if err != nil {
		return nil, err
	}

	// 3. fazer o process do amount(valor $) do invoice
	if err = invoice.Process(); err != nil {
		return nil, err
	}

	// 4. se status é aprovado atualizar o balance da ACCOUNT
	if invoice.Status == domain.StatusApproved {
		_, err := s.accountService.UpdateBalance(input.APIKey, invoice.Amount)
		if err != nil {
			return nil, err
		}
	}

	// 5. salvar o invoice no banco (tabela de faturas)
	if err = s.invoiceRepository.Save(invoice); err != nil {
		return nil, err
	}

	// 6. devolver como dto
	return dto.FromInvoice(invoice), nil

}

func (s *InvoiceService) GetById(id, apiKey string) (*dto.InvoiceOutput, error) {
	invoice, err := s.invoiceRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	accountOutput, err := s.accountService.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}
	
	if invoice.AccountId != accountOutput.ID {
		return nil, domain.ErrUnauthorizedAccess
	}

	return dto.FromInvoice(invoice), nil
}

func (s *InvoiceService) ListByAccountApiKey(apiKey string) ([]*dto.InvoiceOutput, error) {
	accountOutput, err := s.accountService.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	return s.ListByAccount(accountOutput.ID)
}

// func auxiliar para ListByAccountApiKey
func (s *InvoiceService) ListByAccount(accountId string) ([]*dto.InvoiceOutput, error) {
	invoices, err := s.invoiceRepository.FindByAccountId(accountId)
	if err != nil {
		return nil, err
	}

	invoiceOutput := make([]*dto.InvoiceOutput, len(invoices))
	for i, invoice := range invoices {
		invoiceOutput[i] = dto.FromInvoice(invoice)
	}

	return invoiceOutput, nil
}