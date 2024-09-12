package services

import (
	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
)

type PostbackService struct {
	postbackRepo repository.IPostbackRepository
}

type IPostbackService interface {
	IsPostback(string) bool
	Get(string) (*entity.Postback, error)
}

func NewPostbackService(postbackRepo repository.IPostbackRepository) *PostbackService {
	return &PostbackService{
		postbackRepo: postbackRepo,
	}
}

func (s *PostbackService) IsPostback(subkey string) bool {
	count, _ := s.postbackRepo.CountBySubkey(subkey)
	return count > 0
}

func (s *PostbackService) Get(subkey string) (*entity.Postback, error) {
	return s.postbackRepo.GetBySubKey(subkey)
}
