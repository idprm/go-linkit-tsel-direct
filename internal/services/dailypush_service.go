package services

import (
	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
)

type DailypushService struct {
	dailypushRepo repository.IDailypushRepository
}

type IDailypushService interface {
	Save(*entity.Dailypush) error
	Update(*entity.Dailypush) error
}

func NewDailypushService(dailypushRepo repository.IDailypushRepository) *DailypushService {
	return &DailypushService{
		dailypushRepo: dailypushRepo,
	}
}

func (s *DailypushService) Save(t *entity.Dailypush) error {
	err := s.dailypushRepo.Save(t)
	if err != nil {
		return err
	}
	return nil
}

func (s *DailypushService) Update(t *entity.Dailypush) error {
	data := &entity.Dailypush{
		ServiceID: t.ServiceID,
		Msisdn:    t.Msisdn,
		Subject:   t.Subject,
		IsCharge:  false,
	}
	errDelete := s.dailypushRepo.Delete(data)
	if errDelete != nil {
		return errDelete
	}

	errSave := s.dailypushRepo.Save(t)
	if errSave != nil {
		return errSave
	}
	return nil
}
