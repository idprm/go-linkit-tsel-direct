package services

import (
	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
)

type ScheduleService struct {
	scheduleRepo repository.IScheduleRepository
}

type IScheduleService interface {
	GetLocked(string, string) bool
	GetUnlocked(string, string) bool
	UpdateSchedule(bool, string)
}

func NewScheduleService(scheduleRepo repository.IScheduleRepository) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
	}
}

func (s *ScheduleService) GetLocked(name string, hour string) bool {
	count, _ := s.scheduleRepo.CountLocked(name, hour)
	return count > 0
}

func (s *ScheduleService) GetUnlocked(name string, hour string) bool {
	count, _ := s.scheduleRepo.CountUnlocked(name, hour)
	return count > 0
}

func (s *ScheduleService) UpdateSchedule(unlocked bool, name string) {
	s.scheduleRepo.Update(
		&entity.Schedule{
			IsUnlocked: unlocked,
			Name:       name,
		},
	)
}
