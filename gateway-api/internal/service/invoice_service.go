package service

import (
	"github.com/j-ordep/gateway/go-gateway/internal/domain"
	"github.com/j-ordep/gateway/go-gateway/internal/repository"
)

type InvoiceService struct {
	invoiceRepository *repository.InvoiceRepository
	accountService    *AccountService
}

func NewInvoiceService(invoiceRepository repository.InvoiceRepository, accountService AccountService) *InvoiceService {
	return &InvoiceService{
		invoiceRepository: &invoiceRepository,
		accountService:    &accountService,
	}
}

func Create() {

}

func (s *InvoiceService) UpdateStatus(newStatus domain.Status) {
	var invoice *domain.Invoice

	invoice.UpdateStatus(newStatus)

	s.invoiceRepository.UpdateStatus(invoice)
}
