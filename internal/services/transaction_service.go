package services

import (
	"log"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
)

type TransactionService struct {
	transactionRepo repository.ITransactionRepository
}

type ITransactionService interface {
	SaveTransaction(*entity.Transaction) error
	UpdateTransaction(*entity.Transaction) error
	GroupByStatusTransaction() (*[]entity.Transaction, error)
	GroupByStatusDetailTransaction() (*[]entity.Transaction, error)
	GroupByAdnetTransaction() (*[]entity.Transaction, error)
	SelectTransactionToCSV() (*[]entity.TransactionToCSV, error)
}

func NewTransactionService(transactionRepo repository.ITransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) SaveTransaction(t *entity.Transaction) error {
	err := s.transactionRepo.Save(t)
	if err != nil {
		return err
	}
	return nil
}

func (s *TransactionService) UpdateTransaction(t *entity.Transaction) error {
	data := &entity.Transaction{
		ServiceID: t.ServiceID,
		Msisdn:    t.Msisdn,
		Subject:   t.Subject,
		Status:    "FAILED",
	}
	errDelete := s.transactionRepo.Delete(data)
	if errDelete != nil {
		return errDelete
	}

	errSave := s.transactionRepo.Save(t)
	if errSave != nil {
		return errSave
	}
	return nil
}

func (s *TransactionService) GroupByStatusTransaction() (*[]entity.Transaction, error) {
	result, err := s.transactionRepo.SelectByStatus()
	if err != nil {
		return nil, err
	}
	var transactions []entity.Transaction
	if len(*result) > 0 {
		for _, a := range *result {
			transaction := entity.Transaction{
				CreatedAt: a.CreatedAt,
				ServiceID: a.ServiceID,
				Subject:   a.Subject,
				Status:    a.Status,
				Msisdn:    a.Msisdn,
				PIN:       a.GetAmountWithSeparator(),
			}
			transactions = append(transactions, transaction)
		}
	}
	return &transactions, nil
}

func (s *TransactionService) GroupByStatusDetailTransaction() (*[]entity.Transaction, error) {
	trans, err := s.transactionRepo.SelectByStatusDetail()
	if err != nil {
		return nil, err
	}
	return trans, nil
}

func (s *TransactionService) GroupByAdnetTransaction() (*[]entity.Transaction, error) {
	trans, err := s.transactionRepo.SelectByAdnet()
	if err != nil {
		return nil, err
	}
	return trans, nil
}

func (s *TransactionService) SelectTransactionToCSV() (*[]entity.TransactionToCSV, error) {
	result, err := s.transactionRepo.SelectTransactionToCSV()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var trans []entity.TransactionToCSV
	if len(*result) > 0 {
		for _, a := range *result {
			tr := entity.TransactionToCSV{
				Country:          a.Country,
				Operator:         a.Operator,
				Source:           a.Source,
				Msisdn:           a.Msisdn,
				Event:            a.Event,
				EventDate:        a.EventDate,
				Cycle:            a.Cycle,
				Revenue:          a.Revenue,
				ChargeDate:       a.ChargeDate,
				Currency:         a.Currency,
				Publisher:        a.Publisher,
				Handset:          a.Handset,
				Browser:          a.Browser,
				TrxId:            a.TrxId,
				TelcoApiUrl:      a.TelcoApiUrl,
				TelcoApiResponse: a.TelcoApiResponse,
				SmsContent:       a.SmsContent,
				StatusSms:        a.StatusSms,
			}

			tr.SetService(a.Service, a.CampSubKeyword)
			tr.SetEventDate(a.EventDate.String)
			tr.SetChargeDate(a.ChargeDate.String)
			tr.SetCycle(a.Cycle)
			trans = append(trans, tr)
		}
	}
	return &trans, nil
}
