package service

import (
	"context"

	"github.com/j-ordep/gateway/go-gateway/internal/domain"
	"github.com/j-ordep/gateway/go-gateway/internal/domain/events"
	"github.com/j-ordep/gateway/go-gateway/internal/dto"
)

type InvoiceService struct {
	invoiceRepository domain.InvoiceRepository
	accountService    AccountService
	kafkaProducer     KafkaProducerInterface
}

func NewInvoiceService(invoiceRepository domain.InvoiceRepository, accountService AccountService, kafkaProducer KafkaProducerInterface) *InvoiceService {
	return &InvoiceService{
		invoiceRepository: invoiceRepository,
		accountService:    accountService,
		kafkaProducer: kafkaProducer,
	}
}

func (s *InvoiceService) Create(input dto.CreateInvoiceInput) (*dto.InvoiceOutput, error) {
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

	if invoice.Status == domain.StatusPending {
		// Criar e publicar evento de transação pendente
		pendingTransaction := events.NewPendingTransaction(invoice.AccountID, invoice.ID, invoice.Amount)

		if err := s.kafkaProducer.SendingPendingTransaction(context.Background(), *pendingTransaction); err != nil {
			return nil, err
		}
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
	invoice, err := s.invoiceRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	accountOutput, err := s.accountService.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	if invoice.AccountID != accountOutput.ID {
		return nil, domain.ErrUnauthorizedAccess
	}

	return dto.FromInvoice(invoice), nil
}

func (s *InvoiceService) ListByAccountApiKey(apiKey string) ([]*dto.InvoiceOutput, error) {
	accountOutput, err := s.accountService.FindByAPIKey(apiKey)
	if err != nil {
		return nil, err
	}

	return s.ListByAccountId(accountOutput.ID)
}

// func auxiliar para ListByAccountApiKey
func (s *InvoiceService) ListByAccountId(accountId string) ([]*dto.InvoiceOutput, error) {
	invoices, err := s.invoiceRepository.FindByAccountID(accountId)
	if err != nil {
		return nil, err
	}

	invoiceOutput := make([]*dto.InvoiceOutput, len(invoices))
	for i, invoice := range invoices {
		invoiceOutput[i] = dto.FromInvoice(invoice)
	}

	return invoiceOutput, nil
}

// ProcessTransactionResult processa o resultado de uma transação após análise de fraude
func (s *InvoiceService) ProcessTransactionResult(invoiceID string, status domain.Status) error {
	invoice, err := s.invoiceRepository.FindByID(invoiceID)
	if err != nil {
		return err
	}

	if err := invoice.UpdateStatus(status); err != nil {
		return err
	}

	if err := s.invoiceRepository.UpdateStatus(invoice); err != nil {
		return err
	}

	if status == domain.StatusApproved {
		account, err := s.accountService.FindByID(invoice.AccountID)
		if err != nil {
			return err
		}

		if _, err := s.accountService.UpdateBalance(account.APIKey, invoice.Amount); err != nil {
			return err
		}
	}

	return nil
}
