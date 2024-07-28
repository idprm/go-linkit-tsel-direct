package services

import (
	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
)

type ContentService struct {
	contentRepo repository.IContentRepository
}

type IContentService interface {
	GetContent(int, string, string) (*entity.Content, error)
}

func NewContentService(contentRepo repository.IContentRepository) *ContentService {
	return &ContentService{
		contentRepo: contentRepo,
	}
}

func (s *ContentService) GetContent(serviceId int, name string, pin string) (*entity.Content, error) {
	result, err := s.contentRepo.Get(serviceId, name)
	if err != nil {
		return nil, err
	}

	var content entity.Content

	if result != nil {
		content = entity.Content{
			Value: result.Value,
			Tid:   result.Tid,
		}
		content.SetPIN(pin)
	}
	return &content, nil
}

func (s *ContentService) GetContentCustom(serviceId int, name, pin, url string) (*entity.Content, error) {
	result, err := s.contentRepo.Get(serviceId, name)
	if err != nil {
		return nil, err
	}

	var content entity.Content

	if result != nil {
		content = entity.Content{
			Value: result.Value,
			Tid:   result.Tid,
		}
		content.SetPIN(pin)
		content.SetLinkPortalMainPlus("https://mindplus.store/linkit360/login")
	}
	return &content, nil
}
