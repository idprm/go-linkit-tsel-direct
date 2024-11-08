package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/idprm/go-linkit-tsel/internal/providers/telco"
	"github.com/idprm/go-linkit-tsel/internal/services"
	"github.com/idprm/go-linkit-tsel/internal/utils/pin_utils"
	"github.com/idprm/go-linkit-tsel/internal/utils/response_utils"
	"github.com/idprm/go-linkit-tsel/internal/utils/uuid_utils"
	"github.com/wiliehidayat87/rmqp"
)

type RetryHandler struct {
	rmq                 rmqp.AMQP
	logger              *logger.Logger
	sub                 *entity.Subscription
	serviceService      services.IServiceService
	contentService      services.IContentService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
	postbackService     services.IPostbackService
	trafficService      services.ITrafficService
	dailypushService    services.IDailypushService
}

func NewRetryHandler(
	rmq rmqp.AMQP,
	logger *logger.Logger,
	sub *entity.Subscription,
	serviceService services.IServiceService,
	contentService services.IContentService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
	postbackService services.IPostbackService,
	trafficService services.ITrafficService,
	dailypushService services.IDailypushService,
) *RetryHandler {
	return &RetryHandler{
		rmq:                 rmq,
		logger:              logger,
		sub:                 sub,
		serviceService:      serviceService,
		contentService:      contentService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
		postbackService:     postbackService,
		trafficService:      trafficService,
		dailypushService:    dailypushService,
	}
}

func (h *RetryHandler) Firstpush() {
	// check if active sub
	if h.subscriptionService.GetActiveSubscription(h.sub.GetServiceId(), h.sub.GetMsisdn()) {

		/**
		 * Generate PIN (portal) before MT sent
		 */
		var pin string
		if h.sub.IsLatestPIN() {
			pin = h.sub.GetLatestPIN()
		} else {
			pin = h.getLatestMsisdn()
		}

		service, _ := h.serviceService.GetServiceId(h.sub.GetServiceId())
		content, _ := h.getContentFirstpush(h.sub.GetServiceId(), pin)
		smsMT := telco.NewTelco(h.logger, h.sub, service, content)
		resp, err := smsMT.SMSbyParam()
		if err != nil {
			log.Println(err)
		}
		trxId := uuid_utils.GenerateTrxId()

		if response_utils.IsSuccess(string(resp)) {

			subSuccess := &entity.Subscription{
				ServiceID:            h.sub.GetServiceId(),
				Msisdn:               h.sub.GetMsisdn(),
				LatestTrxId:          trxId,
				LatestStatus:         STATUS_SUCCESS,
				LatestSubject:        SUBJECT_FIRSTPUSH,
				LatestPIN:            pin,
				Amount:               service.GetPrice(),
				RenewalAt:            time.Now().AddDate(0, 0, service.GetRenewalDay()),
				ChargeAt:             time.Now(),
				Success:              1,
				Failed:               1,
				TotalFirstpush:       1,
				TotalAmountFirstpush: service.GetPrice(),
				IsRetry:              false,
				LatestPayload:        string(resp),
			}

			h.subscriptionService.UpdateSuccessRetry(subSuccess)

			// if first_success_at is null
			if h.subscriptionService.IsFirstSuccess(h.sub.GetServiceId(), h.sub.GetMsisdn()) {
				h.subscriptionService.UpdateFirstSuccess(
					&entity.Subscription{
						ServiceID:      h.sub.GetServiceId(),
						Msisdn:         h.sub.GetMsisdn(),
						FirstSuccessAt: time.Now(),
					},
				)
			}

			transSuccess := &entity.Transaction{
				TxID:           trxId,
				ServiceID:      h.sub.GetServiceId(),
				Msisdn:         h.sub.GetMsisdn(),
				Channel:        h.sub.GetChannel(),
				Keyword:        h.sub.GetLatestKeyword(),
				Amount:         service.GetPrice(),
				PIN:            pin,
				Status:         STATUS_SUCCESS,
				StatusCode:     string(resp),
				StatusDetail:   response_utils.ParseStatusCode(string(resp)),
				Subject:        SUBJECT_FIRSTPUSH,
				Payload:        string(resp),
				CampKeyword:    h.sub.GetCampKeyword(),
				CampSubKeyword: h.sub.GetCampSubKeyword(),
				IpAddress:      h.sub.GetIpAddress(),
			}

			h.transactionService.UpdateTransaction(transSuccess)

			// update traffics_mo if success charge
			h.trafficService.UpdateMOCharge(
				&entity.TrafficMO{
					ServiceID: h.sub.GetServiceId(),
					Msisdn:    h.sub.GetMsisdn(),
					IsCharge:  true,
				},
			)

			// insert to rabbitmq
			jsonDataNotif, _ := json.Marshal(
				&entity.ReqNotifParams{
					Service:      service,
					Subscription: h.sub,
					Action:       "SUB",
					Pin:          pin,
				},
			)
			h.rmq.IntegratePublish(
				RMQ_NOTIF_EXCHANGE,
				RMQ_NOTIF_QUEUE,
				RMQ_DATA_TYPE,
				"",
				string(jsonDataNotif),
			)

			pb := &entity.ReqPostbackParams{
				Subscription: &entity.Subscription{
					LatestTrxId:    trxId,
					ServiceID:      h.sub.GetServiceId(),
					Msisdn:         h.sub.GetMsisdn(),
					LatestKeyword:  h.sub.GetLatestKeyword(),
					LatestSubject:  SUBJECT_FIRSTPUSH,
					LatestPayload:  string(resp),
					CampKeyword:    h.sub.GetCampKeyword(),
					CampSubKeyword: h.sub.GetCampSubKeyword(),
				},
				Service:   service,
				Postback:  &entity.Postback{},
				Action:    "MT_FIRSTPUSH",
				Status:    STATUS_SUCCESS,
				AffSub:    h.sub.GetAffSub(),
				IsSuccess: true,
			}

			if h.postbackService.IsPostback(h.sub.GetCampSubKeyword()) {
				postback, err := h.postbackService.Get(h.sub.GetCampSubKeyword())
				if err != nil {
					log.Println(err.Error())
				}
				pb.Postback = postback
			}

			jsonDataPostback, _ := json.Marshal(pb)
			h.rmq.IntegratePublish(
				RMQ_POSTBACK_MT_EXCHANGE,
				RMQ_POSTBACK_MT_QUEUE,
				RMQ_DATA_TYPE,
				"",
				string(jsonDataPostback),
			)

		} else {
			/**
			* insuff action
			**/
			h.insuff(string(resp))
		}
	}

}

func (h *RetryHandler) Dailypush() {
	// check if active sub
	if h.subscriptionService.GetActiveSubscription(h.sub.GetServiceId(), h.sub.GetMsisdn()) {
		service, _ := h.serviceService.GetServiceId(h.sub.GetServiceId())
		/**
		 * Generate PIN (portal) before MT sent
		 */
		var pin string
		if h.sub.IsLatestPIN() {
			pin = h.sub.GetLatestPIN()
		} else {
			pin = h.getLatestMsisdn()
		}

		content, _ := h.getContentRenewal(h.sub.GetServiceId(), pin)
		smsMT := telco.NewTelco(h.logger, h.sub, service, content)
		resp, err := smsMT.SMSbyParam()
		if err != nil {
			log.Println(err)
		}
		trxId := uuid_utils.GenerateTrxId()

		if response_utils.IsSuccess(string(resp)) {

			subSuccess := &entity.Subscription{
				ServiceID:          h.sub.GetServiceId(),
				Msisdn:             h.sub.GetMsisdn(),
				LatestTrxId:        trxId,
				LatestStatus:       STATUS_SUCCESS,
				LatestSubject:      SUBJECT_RENEWAL,
				LatestPIN:          pin,
				Amount:             service.GetPrice(),
				RenewalAt:          time.Now().AddDate(0, 0, service.GetRenewalDay()),
				ChargeAt:           time.Now(),
				Success:            1,
				Failed:             1,
				TotalRenewal:       1,
				TotalAmountRenewal: service.GetPrice(),
				IsRetry:            false,
				LatestPayload:      string(resp),
			}

			h.subscriptionService.UpdateSuccessRetry(subSuccess)

			// if first_success_at is null
			if h.subscriptionService.IsFirstSuccess(h.sub.GetServiceId(), h.sub.GetMsisdn()) {
				h.subscriptionService.UpdateFirstSuccess(
					&entity.Subscription{
						ServiceID:      h.sub.GetServiceId(),
						Msisdn:         h.sub.GetMsisdn(),
						FirstSuccessAt: time.Now(),
					},
				)
			}

			transSuccess := &entity.Transaction{
				TxID:           trxId,
				ServiceID:      h.sub.GetServiceId(),
				Msisdn:         h.sub.GetMsisdn(),
				Channel:        h.sub.GetChannel(),
				Adnet:          "",
				Keyword:        h.sub.GetLatestKeyword(),
				Amount:         service.GetPrice(),
				PIN:            pin,
				Status:         STATUS_SUCCESS,
				StatusCode:     string(resp),
				StatusDetail:   response_utils.ParseStatusCode(string(resp)),
				Subject:        SUBJECT_RENEWAL,
				Payload:        string(resp),
				CampKeyword:    h.sub.GetCampKeyword(),
				CampSubKeyword: h.sub.GetCampSubKeyword(),
				IpAddress:      h.sub.GetIpAddress(),
			}

			h.transactionService.UpdateTransaction(transSuccess)

			// insert to rabbitmq
			jsonDataNotif, _ := json.Marshal(
				&entity.ReqNotifParams{
					Service:      service,
					Subscription: h.sub,
					Action:       "RENEWAL",
					Pin:          pin,
				},
			)
			h.rmq.IntegratePublish(
				RMQ_NOTIF_EXCHANGE,
				RMQ_NOTIF_QUEUE,
				RMQ_DATA_TYPE,
				"",
				string(jsonDataNotif),
			)

			pb := &entity.ReqPostbackParams{
				Subscription: &entity.Subscription{
					LatestTrxId:    trxId,
					ServiceID:      h.sub.GetServiceId(),
					Msisdn:         h.sub.GetMsisdn(),
					LatestKeyword:  h.sub.GetLatestKeyword(),
					LatestSubject:  SUBJECT_RETRY,
					LatestPayload:  string(resp),
					CampKeyword:    h.sub.GetCampKeyword(),
					CampSubKeyword: h.sub.GetCampSubKeyword(),
				},
				Service:   service,
				Postback:  &entity.Postback{},
				Action:    "MT_DAILYPUSH",
				Status:    STATUS_SUCCESS,
				AffSub:    h.sub.GetAffSub(),
				IsSuccess: true,
			}

			if h.postbackService.IsPostback(h.sub.GetCampSubKeyword()) {
				postback, err := h.postbackService.Get(h.sub.GetCampSubKeyword())
				if err != nil {
					log.Println(err.Error())
				}
				pb.Postback = postback
			}

			// insert to rabbitmq
			jsonDataPostback, _ := json.Marshal(pb)

			h.rmq.IntegratePublish(
				RMQ_POSTBACK_MT_EXCHANGE,
				RMQ_POSTBACK_MT_QUEUE,
				RMQ_DATA_TYPE,
				"",
				string(jsonDataPostback),
			)

			// insert to rabbitmq
			jsonDataDP, _ := json.Marshal(
				&entity.DailypushBodyRequest{
					TxId:           trxId,
					SubscriptionId: h.sub.GetId(),
					ServiceId:      h.sub.GetServiceId(),
					Msisdn:         h.sub.GetMsisdn(),
					Channel:        h.sub.GetChannel(),
					CampKeyword:    h.sub.GetCampKeyword(),
					CampSubKeyword: h.sub.GetCampSubKeyword(),
					Adnet:          h.sub.GetAdnet(),
					PubID:          h.sub.GetPubId(),
					AffSub:         h.sub.GetAffSub(),
					Subject:        SUBJECT_DAILYPUSH,
					StatusCode:     string(resp),
					StatusDetail:   response_utils.ParseStatusCode(string(resp)),
					IsCharge:       true,
					IpAddress:      h.sub.GetIpAddress(),
					Action:         SUBJECT_RETRY,
				},
			)

			h.rmq.IntegratePublish(
				RMQ_DAILYPUSH_EXCHANGE,
				RMQ_DAILYPUSH_QUEUE,
				RMQ_DATA_TYPE,
				"",
				string(jsonDataDP),
			)
		} else {
			/**
			* insuff action
			**/
			h.insuff(string(resp))
		}
	}
}

func (h *RetryHandler) insuff(resp string) {
	if response_utils.IsInsuff(resp) {
		h.subscriptionService.UpdateLatestPayload(
			&entity.Subscription{
				ServiceID:     h.sub.GetServiceId(),
				Msisdn:        h.sub.GetMsisdn(),
				LatestPayload: resp,
			},
		)
	}
}

func (h *RetryHandler) getContentFirstpush(serviceId int, pin string) (*entity.Content, error) {
	return h.contentService.GetContent(serviceId, MT_FIRSTPUSH, pin)
}

func (h *RetryHandler) getContentRenewal(serviceId int, pin string) (*entity.Content, error) {
	return h.contentService.GetContent(serviceId, MT_RENEWAL, pin)
}

func (h *RetryHandler) getLatestMsisdn() string {
	return pin_utils.GetLatestMsisdn(h.sub.Msisdn, 8)
}
