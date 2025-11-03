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

func (s *InvoiceService) Create(input *dto.CreateInvoiceInput) (*dto.InvoiceOutput, error) {
	
	// 1. achar o dono da fatura (accountId)
	accountOutput, err := s.accountService.FindByAPIKey(input.APIKey)
	if err != nil {
		return nil, err
	}

	invoice, err := dto.ToInvoice(input, accountOutput.ID)
	if err != nil {
		return nil, err
	}

	if err = invoice.Process(); err != nil {
		return nil, err
	}

	if invoice.Status == domain.StatusApproved {
		_, err := s.accountService.UpdateBalance(input.APIKey, invoice.Amount)
		if err != nil {
			return nil, err
		}
	}

	if err = s.invoiceRepository.Save(invoice); err != nil {
		return nil, err
	}

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

func (s *InvoiceService) ListByAccountByApiKey(apiKey string) ([]*dto.InvoiceOutput, error) {
	//50min
}

func (s *InvoiceService) UpdateStatus(newStatus domain.Status) (*dto.InvoiceOutput, error) {
	var invoice *domain.Invoice

	err := invoice.UpdateStatus(newStatus)
	if err != nil {
		return nil, err
	}

	err = s.invoiceRepository.UpdateStatus(invoice)
	if err != nil {
		return nil, err
	}

	return dto.FromInvoice(invoice), nil
}
