package services

import (
	"log"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/repository"
)

type SubscriptionService struct {
	subscriptionRepo repository.ISubscriptionRepository
}

type ISubscriptionService interface {
	GetActiveSubscription(int, string) bool
	GetSubscription(int, string) bool
	GetPinSubscription(int) bool
	GetPinActiveSub(string, string) bool
	IsFirstSuccess(int, string) bool
	SelectSubscription(int, string) (*entity.Subscription, error)
	SelectAgeDay(int, string) (int, error)
	SaveSubscription(*entity.Subscription) error
	UpdateSuccess(*entity.Subscription) error
	UpdateFailed(*entity.Subscription) error
	UpdateLatest(*entity.Subscription) error
	UpdateEnable(*entity.Subscription) error
	UpdateDisable(*entity.Subscription) error
	UpdateConfirm(*entity.Subscription) error
	UpdatePurge(*entity.Subscription) error
	UpdateLatestPayload(*entity.Subscription) error
	UpdatePin(*entity.Subscription) error
	UpdateCampByToken(sub *entity.Subscription) error
	UpdateSuccessRetry(*entity.Subscription) error
	UpdateFirstSuccess(*entity.Subscription) error
	UpdateTotalSub(*entity.Subscription) error
	UpdateTotalUnSub(*entity.Subscription) error
	ReminderSubscription() *[]entity.Subscription
	RenewalSubscription() *[]entity.Subscription
	RetryFpSubscription() *[]entity.Subscription
	RetryDpSubscription() *[]entity.Subscription
	RetryInsuffSubscription() *[]entity.Subscription
	TrialSubscription() *[]entity.Subscription
	EmptyCampSubscription() *[]entity.Subscription
	AveragePerUser(string, string, string, string) (*[]entity.AveragePerUserResponse, error)
	SelectSubcriptionToCSV() (*[]entity.SubscriptionToCSV, error)
	SelectSubcriptionPurge() *[]entity.Subscription
}

func NewSubscriptionService(subscriptionRepo repository.ISubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *SubscriptionService) GetActiveSubscription(serviceId int, msisdn string) bool {
	count, _ := s.subscriptionRepo.CountActive(serviceId, msisdn)
	return count > 0
}

func (s *SubscriptionService) GetSubscription(serviceId int, msisdn string) bool {
	count, _ := s.subscriptionRepo.Count(serviceId, msisdn)
	return count > 0
}

func (s *SubscriptionService) GetPinSubscription(pin int) bool {
	count, err := s.subscriptionRepo.CountPin(pin)
	if err != nil {
		log.Println(err)
	}
	return count > 0
}

func (s *SubscriptionService) GetPinActiveSub(category, pin string) bool {
	count, err := s.subscriptionRepo.CountPinActive(category, pin)
	if err != nil {
		log.Println(err)
	}
	return count > 0
}

func (s *SubscriptionService) IsFirstSuccess(serviceId int, msisdn string) bool {
	count, err := s.subscriptionRepo.CountFirstSuccess(serviceId, msisdn)
	if err != nil {
		log.Println(err)
	}
	return count > 0
}

func (s *SubscriptionService) SelectSubscription(serviceId int, msisdn string) (*entity.Subscription, error) {
	sub, err := s.subscriptionRepo.Get(serviceId, msisdn)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) SelectAgeDay(serviceId int, msisdn string) (int, error) {
	count, err := s.subscriptionRepo.GetAgeDay(serviceId, msisdn)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *SubscriptionService) SaveSubscription(sub *entity.Subscription) error {
	err := s.subscriptionRepo.Save(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateSuccess(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateSuccess(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateFailed(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateFailed(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateLatest(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateLatest(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateEnable(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateEnable(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateDisable(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateDisable(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateConfirm(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateConfirm(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdatePurge(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdatePurge(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateLatestPayload(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateLatestPayload(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdatePin(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdatePin(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateCampByToken(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateCampByToken(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateSuccessRetry(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateSuccessRetry(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateFirstSuccess(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateFirstSuccess(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateTotalSub(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateTotalSub(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) UpdateTotalUnSub(sub *entity.Subscription) error {
	err := s.subscriptionRepo.UpdateTotalUnSub(sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) ReminderSubscription() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.Reminder()
	if err != nil {
		log.Println(err)
	}
	return subs
}

func (s *SubscriptionService) RenewalSubscription() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.Renewal()
	if err != nil {
		log.Println(err)
	}
	return subs
}

func (s *SubscriptionService) RetryFpSubscription() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.RetryFp()
	if err != nil {
		log.Println(err)
	}
	return subs
}

func (s *SubscriptionService) RetryDpSubscription() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.RetryDp()
	if err != nil {
		log.Println(err)
	}
	return subs
}

func (s *SubscriptionService) RetryInsuffSubscription() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.RetryInsuff()
	if err != nil {
		log.Println(err)
	}
	return subs
}

func (s *SubscriptionService) TrialSubscription() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.Trial()
	if err != nil {
		log.Println(err)
	}
	return subs
}

func (s *SubscriptionService) EmptyCampSubscription() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.EmptyCamp()
	if err != nil {
		log.Println(err)
	}
	return subs
}

func (s *SubscriptionService) AveragePerUser(start, end, renewal, subkey string) (*[]entity.AveragePerUserResponse, error) {
	result, err := s.subscriptionRepo.AveragePerUser(start, end, renewal, subkey)
	if err != nil {
		return nil, err
	}

	var arpus []entity.AveragePerUserResponse
	if len(*result) > 0 {
		for _, a := range *result {
			arpu := entity.AveragePerUserResponse{
				Name:       a.Name,
				Service:    a.Service,
				Adnet:      a.Adnet,
				Subs:       a.Subs,
				SubsActive: a.SubsActive,
			}
			arpu.SetRevenue(a.Revenue)
			arpus = append(arpus, arpu)
		}
	}
	return &arpus, nil
}

func (s *SubscriptionService) SelectSubcriptionToCSV() (*[]entity.SubscriptionToCSV, error) {
	result, err := s.subscriptionRepo.SelectSubcriptionToCSV()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var subs []entity.SubscriptionToCSV
	if len(*result) > 0 {
		for _, a := range *result {
			sub := entity.SubscriptionToCSV{
				Country:         a.Country,
				Operator:        a.Operator,
				Source:          a.Source,
				Msisdn:          a.Msisdn,
				LatestSubject:   a.LatestSubject,
				Cycle:           a.Cycle,
				Adnet:           a.Adnet,
				Revenue:         a.Revenue,
				SubsDate:        a.SubsDate,
				RenewalDate:     a.RenewalDate,
				FreemiumEndDate: a.FreemiumEndDate,
				UnsubsFrom:      a.UnsubsFrom,
				UnsubsDate:      a.UnsubsDate,
				ServicePrice:    a.ServicePrice,
				Currency:        a.Currency,
				ProfileStatus:   a.ProfileStatus,
				Publisher:       a.Publisher,
				Trxid:           a.Trxid,
				Pixel:           a.Pixel,
				Handset:         a.Handset,
				Browser:         a.Browser,
				AttemptCharging: a.AttemptCharging,
				SuccessBilling:  a.SuccessBilling,
			}

			sub.SetService(a.Service, a.CampSubKeyword)
			sub.SetSubsDate(a.SubsDate.String)
			sub.SetRenewalDate(a.RenewalDate.String)
			sub.SetUnsubsDate(a.UnsubsDate.String)
			sub.SetProfileStatus(a.ProfileStatus)
			sub.SetLatestSubject(a.LatestSubject)
			sub.SetAdnet(a.Adnet)

			subs = append(subs, sub)
		}
	}
	return &subs, nil
}

func (s *SubscriptionService) SelectSubcriptionPurge() *[]entity.Subscription {
	subs, err := s.subscriptionRepo.SelectSubcriptionPurge()
	if err != nil {
		log.Println(err)
	}
	return subs
}
