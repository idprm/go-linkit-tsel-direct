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

type RenewalHandler struct {
	rmq                 rmqp.AMQP
	logger              *logger.Logger
	sub                 *entity.Subscription
	serviceService      services.IServiceService
	contentService      services.IContentService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
}

func NewRenewalHandler(
	rmq rmqp.AMQP,
	logger *logger.Logger,
	sub *entity.Subscription,
	serviceService services.IServiceService,
	contentService services.IContentService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
) *RenewalHandler {

	return &RenewalHandler{
		rmq:                 rmq,
		logger:              logger,
		sub:                 sub,
		serviceService:      serviceService,
		contentService:      contentService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
	}
}

func (h *RenewalHandler) Dailypush() {
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

		var status string
		var isSuccess bool

		if response_utils.ParseStatus(string(resp)) {
			subSuccess := &entity.Subscription{
				ServiceID:          h.sub.GetServiceId(),
				Msisdn:             h.sub.GetMsisdn(),
				LatestTrxId:        trxId,
				LatestSubject:      SUBJECT_RENEWAL,
				LatestStatus:       STATUS_SUCCESS,
				LatestPIN:          pin,
				Amount:             service.GetPrice(),
				RenewalAt:          time.Now().AddDate(0, 0, service.GetRenewalDay()),
				ChargeAt:           time.Now(),
				Success:            1,
				IsRetry:            false,
				TotalRenewal:       1,
				TotalAmountRenewal: service.GetPrice(),
				LatestPayload:      string(resp),
			}

			h.subscriptionService.UpdateSuccess(subSuccess)

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

			h.transactionService.SaveTransaction(transSuccess)

			// insert to rabbitmq
			jsonData, _ := json.Marshal(
				&entity.ReqNotifParams{
					Service:      service,
					Subscription: subSuccess,
					Action:       "RENEWAL",
					Pin:          pin,
				},
			)
			h.rmq.IntegratePublish(
				RMQ_NOTIFEXCHANGE,
				RMQ_NOTIFQUEUE,
				RMQ_DATATYPE,
				"",
				string(jsonData),
			)

			status = STATUS_SUCCESS
			isSuccess = true

		} else {
			subFailed := &entity.Subscription{
				ServiceID:     h.sub.GetServiceId(),
				Msisdn:        h.sub.GetMsisdn(),
				LatestTrxId:   trxId,
				LatestSubject: SUBJECT_RENEWAL,
				LatestStatus:  STATUS_FAILED,
				RenewalAt:     time.Now().AddDate(0, 0, 1),
				RetryAt:       time.Now(),
				Failed:        1,
				IsRetry:       true,
				LatestPayload: string(resp),
			}
			h.subscriptionService.UpdateFailed(subFailed)

			transFailed := &entity.Transaction{
				TxID:           trxId,
				ServiceID:      h.sub.GetServiceId(),
				Msisdn:         h.sub.GetMsisdn(),
				Channel:        h.sub.GetChannel(),
				Adnet:          "",
				Keyword:        h.sub.GetLatestKeyword(),
				Status:         STATUS_FAILED,
				StatusCode:     string(resp),
				StatusDetail:   response_utils.ParseStatusCode(string(resp)),
				Subject:        SUBJECT_RENEWAL,
				Payload:        string(resp),
				CampKeyword:    h.sub.GetCampKeyword(),
				CampSubKeyword: h.sub.GetCampSubKeyword(),
				IpAddress:      h.sub.GetIpAddress(),
			}

			h.transactionService.SaveTransaction(transFailed)

			status = STATUS_FAILED
			isSuccess = false

			/**
			* purge action
			**/
			h.purge(trxId, string(resp))
		}

		// insert to rabbitmq
		jsonData, _ := json.Marshal(
			&entity.ReqPostbackParams{
				Subscription: &entity.Subscription{
					LatestTrxId:    trxId,
					ServiceID:      h.sub.GetServiceId(),
					Msisdn:         h.sub.GetMsisdn(),
					LatestKeyword:  h.sub.GetLatestKeyword(),
					LatestSubject:  SUBJECT_RENEWAL,
					LatestPayload:  string(resp),
					CampKeyword:    h.sub.GetCampKeyword(),
					CampSubKeyword: h.sub.GetCampSubKeyword(),
				},
				Service:   service,
				Action:    "MT_DAILYPUSH",
				Status:    status,
				AffSub:    h.sub.GetAffSub(),
				IsSuccess: isSuccess,
			},
		)
		h.rmq.IntegratePublish(
			RMQ_POSTBACKMTEXCHANGE,
			RMQ_POSTBACKMTQUEUE,
			RMQ_DATATYPE,
			"",
			string(jsonData),
		)

		var subject string
		if response_utils.IsPurge(string(resp)) {
			subject = SUBJECT_PURGE
		} else {
			subject = SUBJECT_DAILYPUSH
		}

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
				AffSub:         h.sub.GetPubId(),
				Subject:        subject,
				StatusCode:     string(resp),
				StatusDetail:   response_utils.ParseStatusCode(string(resp)),
				IsCharge:       isSuccess,
				IpAddress:      h.sub.GetIpAddress(),
				Action:         SUBJECT_RENEWAL,
			},
		)

		h.rmq.IntegratePublish(
			RMQ_DAILYPUSHEXCHANGE,
			RMQ_DAILYPUSHQUEUE,
			RMQ_DATATYPE,
			"",
			string(jsonDataDP),
		)
	}

}

func (h *RenewalHandler) purge(trxId, resp string) {
	if response_utils.IsPurge(resp) {
		h.subscriptionService.UpdatePurge(
			&entity.Subscription{
				ServiceID:   h.sub.GetServiceId(),
				Msisdn:      h.sub.GetMsisdn(),
				PurgeAt:     time.Now(),
				PurgeReason: response_utils.ParseStatusCode(resp),
			},
		)

		// count total unsub
		h.subscriptionService.UpdateTotalUnSub(
			&entity.Subscription{
				ServiceID:  h.sub.GetServiceId(),
				Msisdn:     h.sub.GetMsisdn(),
				TotalUnsub: 1,
			},
		)

		h.transactionService.SaveTransaction(
			&entity.Transaction{
				TxID:           trxId,
				ServiceID:      h.sub.GetServiceId(),
				Msisdn:         h.sub.GetMsisdn(),
				Channel:        h.sub.GetChannel(),
				Adnet:          "",
				Keyword:        h.sub.GetLatestKeyword(),
				Status:         STATUS_SUCCESS,
				StatusCode:     resp,
				StatusDetail:   response_utils.ParseStatusCode(resp),
				Subject:        SUBJECT_PURGE,
				Payload:        resp,
				CampKeyword:    h.sub.GetCampKeyword(),
				CampSubKeyword: h.sub.GetCampSubKeyword(),
				IpAddress:      h.sub.GetIpAddress(),
			},
		)
	}
}

func (h *RenewalHandler) getLatestMsisdn() string {
	return pin_utils.GetLatestMsisdn(h.sub.Msisdn, 8)
}

func (h *RenewalHandler) getContentRenewal(serviceId int, pin string) (*entity.Content, error) {
	return h.contentService.GetContent(serviceId, MT_RENEWAL, pin)
}
