package services

import (
	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
)

type ServiceService struct {
	serviceRepo repository.IServiceRepository
}

type IServiceService interface {
	CheckService(string) bool
	GetServiceId(int) (*entity.Service, error)
	GetServiceByCode(string) (*entity.Service, error)
}

func NewServiceService(serviceRepo repository.IServiceRepository) *ServiceService {
	return &ServiceService{
		serviceRepo: serviceRepo,
	}
}

func (s *ServiceService) CheckService(code string) bool {
	count, _ := s.serviceRepo.CountByCode(code)
	return count > 0
}

func (s *ServiceService) GetServiceId(id int) (*entity.Service, error) {
	result, err := s.serviceRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	var srv entity.Service

	if result != nil {
		srv = entity.Service{
			ID:                  result.ID,
			Category:            result.Category,
			Code:                result.Code,
			Name:                result.Name,
			Package:             result.Package,
			Price:               result.Price,
			ProgramId:           result.ProgramId,
			Sid:                 result.Sid,
			RenewalDay:          result.RenewalDay,
			TrialDay:            result.TrialDay,
			UrlTelco:            result.UrlTelco,
			UrlPortal:           result.UrlPortal,
			UrlCallback:         result.UrlCallback,
			UrlNotifSub:         result.UrlNotifSub,
			UrlNotifUnsub:       result.UrlNotifUnsub,
			UrlNotifRenewal:     result.UrlNotifRenewal,
			UrlPostback:         result.UrlPostback,
			UrlPostbackBillable: result.UrlPostbackBillable,
			UrlPostbackSamMO:    result.UrlPostbackSamMO,
			UrlPostbackSamDN:    result.UrlPostbackSamDN,
			UrlPostbackYlcMO:    result.UrlPostbackYlcMO,
			UrlPostbackYlcMT:    result.UrlPostbackYlcMT,
			UrlPostbackFsMO:     result.UrlPostbackFsMO,
			UrlPostbackFsDN:     result.UrlPostbackFsDN,
		}
	}
	return &srv, nil
}

func (s *ServiceService) GetServiceByCode(code string) (*entity.Service, error) {
	result, err := s.serviceRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	var srv entity.Service

	if result != nil {
		srv = entity.Service{
			ID:                  result.ID,
			Category:            result.Category,
			Code:                result.Code,
			Name:                result.Name,
			Package:             result.Package,
			Price:               result.Price,
			ProgramId:           result.ProgramId,
			Sid:                 result.Sid,
			RenewalDay:          result.RenewalDay,
			TrialDay:            result.TrialDay,
			UrlTelco:            result.UrlTelco,
			UrlPortal:           result.UrlPortal,
			UrlCallback:         result.UrlCallback,
			UrlNotifSub:         result.UrlNotifSub,
			UrlNotifUnsub:       result.UrlNotifUnsub,
			UrlNotifRenewal:     result.UrlNotifRenewal,
			UrlPostback:         result.UrlPostback,
			UrlPostbackBillable: result.UrlPostbackBillable,
			UrlPostbackSamMO:    result.UrlPostbackSamMO,
			UrlPostbackSamDN:    result.UrlPostbackSamDN,
			UrlPostbackYlcMO:    result.UrlPostbackYlcMO,
			UrlPostbackYlcMT:    result.UrlPostbackYlcMT,
			UrlPostbackFsMO:     result.UrlPostbackFsMO,
			UrlPostbackFsDN:     result.UrlPostbackFsDN,
		}
	}
	return &srv, nil
}
