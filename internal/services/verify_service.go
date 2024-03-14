package services

import (
	"strings"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
)

type VerifyService struct {
	verifyRepo repository.IVerifyRepository
}

type IVerifyService interface {
	SetVerify(*entity.Verify) error
	GetVerify(string) (*entity.Verify, error)
}

func NewVerifyService(verifyRepo repository.IVerifyRepository) *VerifyService {
	return &VerifyService{
		verifyRepo: verifyRepo,
	}
}

func (s *VerifyService) SetVerify(t *entity.Verify) error {
	return s.verifyRepo.Set(t)
}

func (s *VerifyService) GetVerify(token string) (*entity.Verify, error) {
	return s.verifyRepo.Get(strings.ToLower(token))
}
