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

	invoiceOutput := dto.FromInvoice(invoice)

	return &invoiceOutput, nil
}

func (s *InvoiceService) FindByAccountId() {
	
}

func (s *InvoiceService) UpdateStatus(newStatus domain.Status) {
	var invoice *domain.Invoice

	invoice.UpdateStatus(newStatus)

	s.invoiceRepository.UpdateStatus(invoice)
}
