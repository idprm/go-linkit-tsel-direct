package services

import (
	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
)

type TrafficService struct {
	trafficRepo repository.ITrafficRepository
}

type ITrafficService interface {
	SaveCampaign(*entity.TrafficCampaign) error
	SaveMO(*entity.TrafficMO) error
}

func NewTrafficService(trafficRepo repository.ITrafficRepository) *TrafficService {
	return &TrafficService{
		trafficRepo: trafficRepo,
	}
}

func (s *TrafficService) SaveCampaign(t *entity.TrafficCampaign) error {
	err := s.trafficRepo.SaveCampaign(t)
	if err != nil {
		return err
	}
	return nil
}

func (s *TrafficService) SaveMO(t *entity.TrafficMO) error {
	err := s.trafficRepo.SaveMO(t)
	if err != nil {
		return err
	}
	return nil
}
