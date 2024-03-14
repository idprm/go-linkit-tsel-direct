package services

import "github.com/idprm/go-linkit-tsel/src/domain/repository"

type BlacklistService struct {
	blacklistRepo repository.IBlacklistRepository
}

type IBlacklistService interface {
	GetBlacklist(msisdn string) bool
}

func NewBlacklistService(blacklistRepo repository.IBlacklistRepository) *BlacklistService {
	return &BlacklistService{
		blacklistRepo: blacklistRepo,
	}
}

func (s *BlacklistService) GetBlacklist(msisdn string) bool {
	count, _ := s.blacklistRepo.Count(msisdn)
	return count > 0
}
