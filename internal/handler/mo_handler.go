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
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
)

type MOHandler struct {
	rmq                 rmqp.AMQP
	logger              *logger.Logger
	blacklistService    services.IBlacklistService
	serviceService      services.IServiceService
	verifyService       services.IVerifyService
	contentService      services.IContentService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
	historyService      services.IHistoryService
	req                 *entity.ReqMOParams
}

func NewMOHandler(
	rmq rmqp.AMQP,
	logger *logger.Logger,
	blacklistService services.IBlacklistService,
	serviceService services.IServiceService,
	verifyService services.IVerifyService,
	contentService services.IContentService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
	historyService services.IHistoryService,
	req *entity.ReqMOParams,
) *MOHandler {
	return &MOHandler{
		rmq:                 rmq,
		logger:              logger,
		blacklistService:    blacklistService,
		serviceService:      serviceService,
		verifyService:       verifyService,
		contentService:      contentService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
		historyService:      historyService,
		req:                 req,
	}
}

func (h *MOHandler) Firstpush() {
	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}

	/**
	 * Generate PIN (portal) before MT sent
	 */
	pin := h.getLatestMsisdn()

	content, err := h.getContentFirstpush(service.GetID(), pin)
	if err != nil {
		log.Println(err)
	}

	channel := response_utils.ParseChannel(h.req.SMS)
	trxId := uuid_utils.GenerateTrxId()

	token := response_utils.ParseToken(h.req.SMS)
	verify, err := h.verifyService.GetVerify(token)
	if err != nil {
		log.Println(err)
	}

	subscription := &entity.Subscription{
		ServiceID:     service.GetID(),
		Category:      service.GetCategory(),
		Msisdn:        h.req.GetMsisdn(),
		LatestTrxId:   trxId,
		LatestKeyword: h.req.GetKeyword(),
		LatestSubject: SUBJECT_FIRSTPUSH,
		Channel:       channel,
		IsActive:      true,
	}

	if verify != nil {
		subscription.Adnet = verify.GetAdnet()
		subscription.PubID = verify.GetPubId()
		subscription.AffSub = verify.GetAffSub()
		subscription.CampKeyword = verify.GetCampKeyword()
		subscription.CampSubKeyword = verify.GetCampSubKeyword()
		subscription.IpAddress = verify.GetIpAddress()

		// insert to rabbitmq
		jsonDataPostback, _ := json.Marshal(
			&entity.ReqPostbackParams{
				Verify:       verify,
				Subscription: subscription,
				Service:      service,
				Action:       "MO",
			},
		)
		h.rmq.IntegratePublish(
			RMQ_POSTBACKMOEXCHANGE,
			RMQ_POSTBACKMOQUEUE,
			RMQ_DATATYPE,
			"",
			string(jsonDataPostback),
		)

	} else {
		subscription.Adnet = ""
		subscription.PubID = ""
		subscription.AffSub = ""
		subscription.CampKeyword = ""
		subscription.CampSubKeyword = ""
		subscription.IpAddress = ""
	}

	if h.IsSub() {
		h.subscriptionService.UpdateEnable(subscription)

	} else {
		h.subscriptionService.SaveSubscription(subscription)
	}

	smsMT := telco.NewTelco(h.logger, subscription, service, content)
	resp, err := smsMT.SMSbyParam()
	if err != nil {
		log.Println(err.Error())
	}

	var status string
	var isSuccess bool

	if response_utils.ParseStatus(string(resp)) {
		subSuccess := &entity.Subscription{
			ServiceID:            service.GetID(),
			Msisdn:               h.req.GetMsisdn(),
			LatestTrxId:          trxId,
			LatestSubject:        SUBJECT_FIRSTPUSH,
			LatestStatus:         STATUS_SUCCESS,
			LatestPIN:            pin,
			Amount:               service.GetPrice(),
			RenewalAt:            time.Now().AddDate(0, 0, service.GetRenewalDay()),
			ChargeAt:             time.Now(),
			Success:              1,
			IsRetry:              false,
			TotalFirstpush:       1,
			TotalAmountFirstpush: service.GetPrice(),
			LatestPayload:        string(resp),
		}

		h.subscriptionService.UpdateSuccess(subSuccess)

		transSuccess := &entity.Transaction{
			TxID:         trxId,
			ServiceID:    service.GetID(),
			Msisdn:       h.req.GetMsisdn(),
			Channel:      channel,
			Keyword:      h.req.GetKeyword(),
			Amount:       service.GetPrice(),
			PIN:          pin,
			Status:       STATUS_SUCCESS,
			StatusCode:   string(resp),
			StatusDetail: response_utils.ParseStatusCode(string(resp)),
			Subject:      SUBJECT_FIRSTPUSH,
			Payload:      string(resp),
		}

		if verify != nil {
			transSuccess.Adnet = verify.GetAdnet()
			transSuccess.PubID = verify.GetPubId()
			transSuccess.AffSub = verify.GetAffSub()
			transSuccess.CampKeyword = verify.GetCampKeyword()
			transSuccess.CampSubKeyword = verify.GetCampSubKeyword()
			transSuccess.IpAddress = verify.GetIpAddress()
		}

		h.transactionService.SaveTransaction(transSuccess)

		historySuccess := &entity.History{
			ServiceID: service.GetID(),
			Msisdn:    h.req.GetMsisdn(),
			Channel:   channel,
			Keyword:   h.req.GetKeyword(),
			Subject:   SUBJECT_FIRSTPUSH,
			Status:    STATUS_SUCCESS,
		}

		if verify != nil {
			historySuccess.Adnet = verify.GetAdnet()
			historySuccess.IpAddress = verify.GetIpAddress()
		}

		h.historyService.SaveHistory(historySuccess)

		// insert to rabbitmq
		jsonDataNotif, _ := json.Marshal(
			&entity.ReqNotifParams{
				Service:      service,
				Subscription: subscription,
				Action:       "SUB",
				Pin:          pin,
			},
		)
		h.rmq.IntegratePublish(
			RMQ_NOTIFEXCHANGE,
			RMQ_NOTIFQUEUE,
			RMQ_DATATYPE,
			"",
			string(jsonDataNotif),
		)

		status = STATUS_SUCCESS
		isSuccess = true

	} else {

		subFailed := &entity.Subscription{
			ServiceID:     service.GetID(),
			Msisdn:        h.req.GetMsisdn(),
			LatestTrxId:   trxId,
			LatestSubject: SUBJECT_FIRSTPUSH,
			LatestStatus:  STATUS_FAILED,
			RenewalAt:     time.Now().AddDate(0, 0, 1),
			RetryAt:       time.Now(),
			Failed:        1,
			IsRetry:       true,
			LatestPayload: string(resp),
		}
		h.subscriptionService.UpdateFailed(subFailed)

		// keep update PIN if failed
		h.subscriptionService.UpdatePin(
			&entity.Subscription{
				ServiceID: service.GetID(),
				Msisdn:    h.req.GetMsisdn(),
				LatestPIN: pin,
			},
		)

		transFailed := &entity.Transaction{
			TxID:         trxId,
			ServiceID:    service.GetID(),
			Msisdn:       h.req.GetMsisdn(),
			Channel:      channel,
			Keyword:      h.req.GetKeyword(),
			Status:       STATUS_FAILED,
			StatusCode:   string(resp),
			StatusDetail: response_utils.ParseStatusCode(string(resp)),
			Subject:      SUBJECT_FIRSTPUSH,
			Payload:      string(resp),
		}

		if verify != nil {
			transFailed.Adnet = verify.GetAdnet()
			transFailed.PubID = verify.GetPubId()
			transFailed.AffSub = verify.GetAffSub()
			transFailed.CampKeyword = verify.GetCampKeyword()
			transFailed.CampSubKeyword = verify.GetCampSubKeyword()
			transFailed.IpAddress = verify.GetIpAddress()
		}
		h.transactionService.SaveTransaction(transFailed)

		historyFailed := &entity.History{
			ServiceID: service.GetID(),
			Msisdn:    h.req.GetMsisdn(),
			Channel:   channel,
			Keyword:   h.req.GetKeyword(),
			Subject:   SUBJECT_FIRSTPUSH,
			Status:    STATUS_FAILED,
		}
		if verify != nil {
			historyFailed.Adnet = verify.GetAdnet()
			historyFailed.IpAddress = verify.GetIpAddress()
		}
		h.historyService.SaveHistory(historyFailed)

		status = STATUS_FAILED
		isSuccess = false
	}

	// postback queue
	if verify != nil {
		// insert to rabbitmq
		jsonDataPostback, _ := json.Marshal(
			&entity.ReqPostbackParams{
				Verify:       verify,
				Subscription: subscription,
				Service:      service,
				Action:       "MT",
				Status:       status,
				IsSuccess:    isSuccess,
			},
		)
		h.rmq.IntegratePublish(
			RMQ_POSTBACKMOEXCHANGE,
			RMQ_POSTBACKMOQUEUE,
			RMQ_DATATYPE,
			"",
			string(jsonDataPostback),
		)
	}
}

func (h *MOHandler) Unsub() {
	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}
	channel := response_utils.ParseChannel(h.req.SMS)
	trxId := uuid_utils.GenerateTrxId()

	subscription := &entity.Subscription{
		ServiceID:     service.GetID(),
		Msisdn:        h.req.GetMsisdn(),
		Channel:       channel,
		LatestTrxId:   trxId,
		LatestKeyword: h.req.GetKeyword(),
		LatestSubject: SUBJECT_UNSUB,
		LatestStatus:  STATUS_SUCCESS,
		UnsubAt:       time.Now(),
		IpAddress:     h.req.GetIpAddress(),
		IsRetry:       false,
		IsActive:      false,
	}
	h.subscriptionService.UpdateDisable(subscription)

	// if unsub, set PIN to 0
	h.subscriptionService.UpdatePin(&entity.Subscription{
		ServiceID: service.GetID(),
		Msisdn:    h.req.GetMsisdn(),
		LatestPIN: "",
	})

	confirm := &entity.Subscription{
		ServiceID: service.GetID(),
		Msisdn:    h.req.GetMsisdn(),
		IsConfirm: false,
	}

	h.subscriptionService.UpdateConfirm(confirm)

	// select data by service_id & msisdn
	sub, _ := h.subscriptionService.SelectSubscription(service.GetID(), h.req.GetMsisdn())

	transaction := &entity.Transaction{
		TxID:         trxId,
		ServiceID:    service.GetID(),
		Msisdn:       h.req.GetMsisdn(),
		Channel:      channel,
		Adnet:        sub.GetAdnet(),
		Keyword:      h.req.GetKeyword(),
		Status:       STATUS_SUCCESS,
		StatusCode:   "-",
		StatusDetail: "-",
		Subject:      SUBJECT_UNSUB,
		Payload:      "-",
	}

	if sub != nil {
		transaction.SetCampKeyword(sub.GetCampKeyword())
		transaction.SetCampSubKeyword(sub.GetCampSubKeyword())
	}

	h.transactionService.SaveTransaction(transaction)

	history := &entity.History{
		ServiceID: service.GetID(),
		Msisdn:    h.req.GetMsisdn(),
		Channel:   channel,
		Adnet:     sub.GetAdnet(),
		Keyword:   h.req.GetKeyword(),
		Subject:   SUBJECT_UNSUB,
		Status:    STATUS_SUCCESS,
		IpAddress: h.req.GetIpAddress(),
	}
	h.historyService.SaveHistory(history)

	// insert to rabbitmq
	jsonDataNotif, _ := json.Marshal(
		&entity.ReqNotifParams{
			Service:      service,
			Subscription: subscription,
			Action:       "UNSUB",
		},
	)

	// insert to rabbitmq
	jsonDataPostback, _ := json.Marshal(
		&entity.ReqPostbackParams{
			Verify:       &entity.Verify{},
			Subscription: sub,
			Service:      service,
			Action:       "MO_UNSUB",
		},
	)

	h.rmq.IntegratePublish(
		RMQ_NOTIFEXCHANGE,
		RMQ_NOTIFQUEUE,
		RMQ_DATATYPE,
		"",
		string(jsonDataNotif),
	)

	h.rmq.IntegratePublish(
		RMQ_POSTBACKMOEXCHANGE,
		RMQ_POSTBACKMOQUEUE,
		RMQ_DATATYPE,
		"",
		string(jsonDataPostback),
	)
}

func (h *MOHandler) Confirm() {
	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}
	subscription := &entity.Subscription{
		ServiceID: service.GetID(),
		Msisdn:    h.req.GetMsisdn(),
		IsConfirm: true,
	}
	h.subscriptionService.UpdateConfirm(subscription)
}

func (h *MOHandler) getService() (*entity.Service, error) {
	keyword := h.req.GetSubKeyword()
	return h.serviceService.GetServiceByCode(keyword)
}

func (h *MOHandler) IsActiveSub() bool {
	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}
	return h.subscriptionService.GetActiveSubscription(service.GetID(), h.req.GetMsisdn())
}

func (h *MOHandler) IsSub() bool {
	service, err := h.getService()
	if err != nil {
		log.Println(err)
	}
	return h.subscriptionService.GetSubscription(service.GetID(), h.req.GetMsisdn())
}

func (h *MOHandler) IsBlacklist() bool {
	return h.blacklistService.GetBlacklist(h.req.GetMsisdn())
}

func (h *MOHandler) Logger(req *entity.ReqMOParams, data string) {
	l := h.logger.Init("mo", true)
	l.WithFields(logrus.Fields{"request": req}).Info(data)
}

func (h *MOHandler) IsService() bool {
	subKeyword := h.req.GetSubKeyword()
	return h.serviceService.CheckService(subKeyword)
}

func (h *MOHandler) getContentFirstpush(serviceId int, pin string) (*entity.Content, error) {
	return h.contentService.GetContent(serviceId, MT_FIRSTPUSH, pin)
}

func (h *MOHandler) getLatestMsisdn() string {
	return pin_utils.GetLatestMsisdn(h.req.Msisdn, 8)
}
