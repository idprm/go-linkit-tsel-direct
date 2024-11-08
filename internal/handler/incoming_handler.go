package handler

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/model"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/idprm/go-linkit-tsel/internal/providers/telco"
	"github.com/idprm/go-linkit-tsel/internal/services"
	"github.com/idprm/go-linkit-tsel/internal/utils"
	"github.com/idprm/go-linkit-tsel/internal/utils/response_utils"
	"github.com/idprm/go-linkit-tsel/internal/utils/uuid_utils"
	"github.com/mileusna/useragent"
	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
)

var (
	APP_URL      string = utils.GetEnv("APP_URL")
	TELCO_SENDER string = utils.GetEnv("TELCO_SENDER")
)

type IncomingHandler struct {
	logger              *logger.Logger
	message             rmqp.AMQP
	serviceService      services.IServiceService
	verifyService       services.IVerifyService
	subscriptionService services.ISubscriptionService
	transactionService  services.ITransactionService
}

func NewIncomingHandler(
	logger *logger.Logger,
	message rmqp.AMQP,
	serviceService services.IServiceService,
	verifyService services.IVerifyService,
	subscriptionService services.ISubscriptionService,
	transactionService services.ITransactionService,
) *IncomingHandler {
	return &IncomingHandler{
		logger:              logger,
		message:             message,
		serviceService:      serviceService,
		verifyService:       verifyService,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
	}
}

const (
	RMQ_DATA_TYPE            string = "application/json"
	RMQ_MO_EXCHANGE          string = "E_MO"
	RMQ_MO_QUEUE             string = "Q_MO"
	RMQ_NOTIF_EXCHANGE       string = "E_NOTIF"
	RMQ_NOTIF_QUEUE          string = "Q_NOTIF"
	RMQ_POSTBACK_MO_EXCHANGE string = "E_POSTBACK_MO"
	RMQ_POSTBACK_MO_QUEUE    string = "Q_POSTBACK_MO"
	RMQ_POSTBACK_MT_EXCHANGE string = "E_POSTBACK_MT"
	RMQ_POSTBACK_MT_QUEUE    string = "Q_POSTBACK_MT"
	RMQ_POSTBACK_FP_EXCHANGE string = "E_POSTBACK_FP"
	RMQ_POSTBACK_FP_QUEUE    string = "Q_POSTBACK_FP"
	RMQ_TRAFFIC_EXCHANGE     string = "E_TRAFFIC"
	RMQ_TRAFFIC_QUEUE        string = "Q_TRAFFIC"
	RMQ_DAILYPUSH_EXCHANGE   string = "E_BQ_DAILYPUSH"
	RMQ_DAILYPUSH_QUEUE      string = "Q_BQ_DAILYPUSH"
	MT_FIRSTPUSH             string = "FIRSTPUSH"
	MT_RENEWAL               string = "RENEWAL"
	MT_UNSUB                 string = "UNSUB"
	STATUS_SUCCESS           string = "SUCCESS"
	STATUS_FAILED            string = "FAILED"
	SUBJECT_FIRSTPUSH        string = "FIRSTPUSH"
	SUBJECT_DAILYPUSH        string = "DAILYPUSH"
	SUBJECT_RENEWAL          string = "RENEWAL"
	SUBJECT_UNSUB            string = "UNSUB"
	SUBJECT_RETRY            string = "RETRY"
	SUBJECT_PURGE            string = "PURGE"
)

var validate = validator.New()

func ValidateStruct(data interface{}) []*entity.ErrorResponse {
	var errors []*entity.ErrorResponse
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element entity.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func (h *IncomingHandler) CloudPlaySubPage(c *fiber.Ctx) error {
	return c.Render("cloudplay/sub", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) GalaysSubPage(c *fiber.Ctx) error {
	return c.Render("galays/sub", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) CloudPlaySub1Page(c *fiber.Ctx) error {
	return c.Render("cloudplay/sub1", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) GalaysSub1Page(c *fiber.Ctx) error {
	return c.Render("galays/sub1", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) CloudPlaySub2Page(c *fiber.Ctx) error {
	return c.Render("cloudplay/sub2", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) CloudPlaySub3Page(c *fiber.Ctx) error {
	return c.Render("cloudplay/sub3", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) CloudPlaySub4Page(c *fiber.Ctx) error {
	return c.Render("cloudplay/sub4", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) CloudPlayUnsubPage(c *fiber.Ctx) error {
	return c.Render("cloudplay/unsub", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) GalaysUnsubPage(c *fiber.Ctx) error {
	return c.Render("galays/unsub", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) CloudPlayTermPage(c *fiber.Ctx) error {
	return c.Render("cloudplay/term", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) CloudPlayCampaign(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService("CLOUDPLAY") {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode("CLOUDPLAY")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        service.GetCode(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     false,
			IsCampTool:     false,
		},
	)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": err,
			},
		)
	}
	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"ip_address":   req.GetIpAddress(),
			"duration":     duration,
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Redirect(redirect, 303)
}

func (h *IncomingHandler) GalaysCampaign(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService("GALAYS") {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode("GALAYS")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        service.GetCode(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     false,
			IsCampTool:     false,
		},
	)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": err,
			},
		)
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   req.GetIpAddress(),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Redirect(redirect, 303)
}

func (h *IncomingHandler) CloudPlayCampaignBillable(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService("CLOUDPLAY") {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode("CLOUDPLAY")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")

		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": "Failed",
			},
		)
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        service.GetCode(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     true,
			IsCampTool:     false,
		},
	)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": err,
			},
		)
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Redirect(redirect, 303)
}

func (h *IncomingHandler) GalaysCampaignBillable(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService("GALAYS") {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode("GALAYS")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        service.GetCode(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     true,
			IsCampTool:     false,
		},
	)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Redirect(redirect, 303)
}

func (h *IncomingHandler) CampaignTool(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()

	req := new(entity.CampaignToolsRequest)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService(req.GetService()) {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode(req.GetService())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"os":           req.GetOS(),
			"browser":      req.GetBrowser(),
			"device":       "",
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        service.GetCode(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    "REG " + req.GetService(),
			CampSubKeyword: req.GetSubKeyword(),
			Browser:        req.GetBrowser(),
			OS:             req.GetOS(),
			Device:         "",
			IpAddress:      req.GetIpAddress(),
			IsBillable:     req.IsBillable(),
			IsCampTool:     true,
		},
	)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    "REG " + req.GetService(),
			CampSubKeyword: req.GetSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        req.GetBrowser(),
			OS:             req.GetOS(),
			Device:         "",
			Referer:        req.GetReferer(),
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   req.GetIpAddress(),
			"os":           req.GetOS(),
			"browser":      req.GetBrowser(),
			"device":       "",
		}).Info("REDIRECT")
	}

	return c.Status(fiber.StatusOK).JSON(
		&entity.CampaignToolsResponse{
			StatusCode: 0,
			Token:      strings.TrimSpace(token),
			UrlPromo:   strings.TrimSpace(redirect),
		},
	)
}

func (h *IncomingHandler) CloudPlaySub1CampaignPage(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService("CLOUDPLAY1") {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode("CLOUDPLAY1")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": "Failed",
			},
		)
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        service.GetCode(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     false,
			IsCampTool:     false,
		},
	)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": err,
			},
		)
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Redirect(redirect, 303)
}

func (h *IncomingHandler) GalaysSub1CampaignPage(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService("GALAYS1") {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode("GALAYS1")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        service.GetCode(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     false,
			IsCampTool:     false,
		},
	)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Redirect(redirect, 303)
}

func (h *IncomingHandler) CloudPlaySub2CampaignPage(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService("CLOUDPLAY2") {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode("CLOUDPLAY2")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": "Failed",
			},
		)
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        "CLOUDPLAY2",
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     false,
			IsCampTool:     false,
		},
	)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Redirect(redirect, 303)
}

func (h *IncomingHandler) CloudPlaySub3CampaignPage(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService("CLOUDPLAY3") {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode("CLOUDPLAY3")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": "Failed",
			},
		)
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        service.GetCode(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     false,
			IsCampTool:     false,
		},
	)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": err,
			},
		)
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Redirect(redirect, 303)
}

func (h *IncomingHandler) CloudPlaySub4CampaignPage(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService("CLOUDPLAY4") {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode("CLOUDPLAY4")
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        service.GetCode(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     false,
			IsCampTool:     false,
		},
	)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": err,
			},
		)
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Redirect(redirect, 303)
}

func (h *IncomingHandler) CampaignToolDynamic(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.CampaignToolsRequest)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService(req.GetDynamic()) {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	var sub *entity.Subscription
	var content *entity.Content

	service, err := h.serviceService.GetServiceByCode(req.GetDynamic())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")

		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        req.GetDynamic(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    "REG " + req.GetDynamic(),
			CampSubKeyword: req.GetSubDynamic(),
			Browser:        req.GetBrowser(),
			OS:             req.GetOS(),
			Device:         "",
			IpAddress:      req.GetIpAddress(),
			IsBillable:     req.IsBillable(),
			IsCampTool:     true,
		},
	)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    "REG " + req.GetDynamic(),
			CampSubKeyword: req.GetSubDynamic(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        req.GetBrowser(),
			OS:             req.GetOS(),
			Device:         "",
			Referer:        req.GetReferer(),
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   req.GetIpAddress(),
			"os":           req.GetOS(),
			"browser":      req.GetBrowser(),
			"device":       "",
		}).Info("REDIRECT")
	}

	return c.Status(fiber.StatusOK).JSON(
		&entity.CampaignToolsResponse{
			StatusCode: 0,
			Token:      strings.TrimSpace(token),
			UrlPromo:   strings.TrimSpace(redirect),
		},
	)
}

func (h *IncomingHandler) CampaignDirect(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	if !h.serviceService.CheckService(strings.ToUpper(c.Params("service"))) {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": "Service Unavailable",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(strings.ToUpper(c.Params("service")))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": "Failed",
			},
		)
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	telco := telco.NewTelco(h.logger, &entity.Subscription{}, service, &entity.Content{})
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   req.GetIpAddress(),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        service.GetCode(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     false,
			IsCampTool:     false,
		},
	)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": err,
			},
		)
	}

	// insert to rabbitmq
	jsonData, _ := json.Marshal(
		&entity.ReqTrafficParams{
			TxId:           trxId,
			ServiceId:      service.GetId(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Adnet:          req.GetAdnet(),
			PubId:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
		},
	)
	h.message.IntegratePublish(
		RMQ_TRAFFIC_EXCHANGE,
		RMQ_TRAFFIC_QUEUE,
		RMQ_DATA_TYPE,
		"",
		string(jsonData),
	)

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   req.GetIpAddress(),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Redirect(redirect, 303)
}

func (h *IncomingHandler) SubPage(c *fiber.Ctx) error {

	req := new(entity.CampaignToolsRequest)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	srv := strings.ToUpper(c.Params("service"))

	if !h.serviceService.CheckService(srv) {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": "Service Unavailable",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(srv)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error":   true,
				"message": "error_internal_server",
			},
		)
	}

	if service.IsCloudplay() {
		return c.Render("cloudplay/sub", fiber.Map{
			"host":         c.BaseURL(),
			"telco_sender": TELCO_SENDER,
			"service_code": srv,
			"gtag":         req.GetGTag(),
			"tiktok":       req.GetTiktok(),
		})
	}

	if service.IsGalays() {
		return c.Render("galays/sub", fiber.Map{
			"host":         c.BaseURL(),
			"telco_sender": TELCO_SENDER,
			"service_code": srv,
			"gtag":         req.GetGTag(),
			"tiktok":       req.GetTiktok(),
		})
	}

	if service.IsMplus() {
		return c.Render("mplus/sub", fiber.Map{
			"host":         c.BaseURL(),
			"telco_sender": TELCO_SENDER,
			"service_code": srv,
			"gtag":         req.GetGTag(),
			"tiktok":       req.GetTiktok(),
		})
	}

	if service.IsGupi() {
		return c.Render("gupi/sub", fiber.Map{
			"host":         c.BaseURL(),
			"telco_sender": TELCO_SENDER,
			"service_code": srv,
			"gtag":         req.GetGTag(),
			"tiktok":       req.GetTiktok(),
		})
	}

	return c.Render("cloudplay/term", fiber.Map{
		"host": c.BaseURL(),
	})
}

func (h *IncomingHandler) TermPage(c *fiber.Ctx) error {
	srv := strings.ToUpper(c.Params("service"))

	if !h.serviceService.CheckService(srv) {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": "Service Unavailable",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(srv)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error":   true,
				"message": "error_internal_server",
			},
		)
	}

	if service.IsGupi() {
		return c.Render("gupi/term", fiber.Map{
			"host":         c.BaseURL(),
			"telco_sender": TELCO_SENDER,
			"service_code": srv,
		})
	}

	return c.Redirect(APP_URL)
}

func (h *IncomingHandler) FaqPage(c *fiber.Ctx) error {
	srv := strings.ToUpper(c.Params("service"))

	if !h.serviceService.CheckService(srv) {
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": "Service Unavailable",
			},
		)
	}

	service, err := h.serviceService.GetServiceByCode(srv)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error":   true,
				"message": "error_internal_server",
			},
		)
	}

	if service.IsMplus() {
		return c.Render("mplus/faq", fiber.Map{
			"host":         c.BaseURL(),
			"telco_sender": TELCO_SENDER,
			"service_code": srv,
		})
	}

	return c.Redirect(APP_URL)
}

func (h *IncomingHandler) OptIn(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	start := time.Now()

	trxId := uuid_utils.GenerateTrxId()
	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.ReqOptInParam)

	err := c.BodyParser(req)
	if err != nil {
		log.Println(err)
	}
	var sub *entity.Subscription
	var content *entity.Content

	if !h.serviceService.CheckService(req.GetService()) {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Service Unavailable",
		})
	}

	service, err := h.serviceService.GetServiceByCode(req.GetService())
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	telco := telco.NewTelco(h.logger, sub, service, content)
	redirect, token, err := telco.WebOptInOTP()
	if err != nil {
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Error("NO_REDIRECT")
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   true,
			"message": "Failed",
		})
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	// insert token & params campaign
	err = h.verifyService.SetVerify(
		&entity.Verify{
			TxId:           trxId,
			Token:          strings.TrimSpace(token),
			Service:        req.GetService(),
			Adnet:          req.GetAdnet(),
			PubID:          req.GetPubId(),
			AffSub:         req.GetAffSub(),
			CampKeyword:    req.GetCampKeyword(),
			CampSubKeyword: req.GetCampSubKeyword(),
			Browser:        ua.Name,
			OS:             ua.OS + " " + ua.OSVersion,
			Device:         ua.Device,
			IpAddress:      req.GetIpAddress(),
			IsBillable:     false,
		},
	)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadGateway).JSON(
			fiber.Map{
				"error":   true,
				"message": err,
			},
		)
	}

	duration := time.Since(start).Milliseconds()
	if token != "" {
		l.WithFields(logrus.Fields{
			"url_campaign": c.OriginalURL(),
			"url_redirect": redirect,
			"duration":     duration,
			"ip_address":   c.Get("X-Forwarded-For"),
			"os":           ua.OS + " " + ua.OSVersion,
			"browser":      ua.Name,
			"device":       ua.Device,
		}).Info("REDIRECT")
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"error":    false,
			"message":  "Success",
			"redirect": redirect,
		},
	)
}

func (h *IncomingHandler) CallbackUrl(c *fiber.Ctx) error {
	l := h.logger.Init("traffic", true)

	userAgent := c.Get("USER-AGENT")
	ua := useragent.Parse(userAgent)

	req := new(entity.SuccessQueryParamsRequest)

	err := c.QueryParser(req)
	if err != nil {
		log.Println(err)
	}

	verify, err := h.verifyService.GetVerify(req.GetToken())
	if err != nil {
		l.WithFields(logrus.Fields{
			"url_success": c.OriginalURL(),
			"error":       err.Error(),
			"ip_address":  c.Get("X-Forwarded-For"),
			"os":          ua.OS + " " + ua.OSVersion,
			"browser":     ua.Name,
			"device":      ua.Device,
		}).Error("PAGE_SUCCESS")

		return c.Render("success", fiber.Map{
			"host": APP_URL,
		})
	}

	l.WithFields(logrus.Fields{
		"url_success": c.OriginalURL(),
		"verify":      verify,
		"ip_address":  c.Get("X-Forwarded-For"),
		"os":          ua.OS + " " + ua.OSVersion,
		"browser":     ua.Name,
		"device":      ua.Device,
	}).Info("PAGE_SUCCESS")

	if !h.serviceService.CheckService(verify.GetService()) {
		return c.Render("success", fiber.Map{
			"host": APP_URL,
		})
	}
	service, _ := h.serviceService.GetServiceByCode(verify.GetService())
	return c.Redirect(service.GetUrlPortal(), 303)

}

func (h *IncomingHandler) MessageOriginated(c *fiber.Ctx) error {
	l := h.logger.Init("mo", true)
	/**
	 * Query Parser
	 */
	req := new(entity.ReqMOParams)

	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			entity.ResponseMO{
				StatusCode: fiber.StatusBadRequest,
				Message:    err.Error(),
			},
		)
	}

	errors := ValidateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.IpAddress = c.Get("Cf-Connecting-Ip")
	} else {
		req.IpAddress = c.Get("X-Forwarded-For")
	}

	json, _ := json.Marshal(req)

	h.message.IntegratePublish(RMQ_MO_EXCHANGE, RMQ_MO_QUEUE, RMQ_DATA_TYPE, "", string(json))

	l.WithFields(logrus.Fields{"request": req}).Info("MO")

	/**
	 * Add New log MO_COMBINE
	 */
	if req.IsREG() {
		token := response_utils.ParseToken(req.SMS)
		verify, err := h.verifyService.GetVerify(token)
		if err != nil {
			l.WithFields(logrus.Fields{
				"request": req,
				"error":   err.Error(),
			}).Error("MO_COMBINE")
		} else {
			l.WithFields(logrus.Fields{
				"request":    req,
				"verify":     verify,
				"ip_address": verify.GetIpAddress(),
				"os":         verify.GetOS(),
				"browser":    verify.GetBrowser(),
				"device":     "",
			}).Info("MO_COMBINE")
		}

	}

	return c.Status(fiber.StatusOK).JSON(entity.ResponseMO{
		StatusCode: fiber.StatusOK,
		Message:    "Successful",
	})
}

func (h *IncomingHandler) Success(c *fiber.Ctx) error {
	return c.Render("success", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) Cancel(c *fiber.Ctx) error {
	return c.Render("cancel", fiber.Map{
		"host": APP_URL,
	})
}

func (h *IncomingHandler) SelectStatus(c *fiber.Ctx) error {
	transactions, err := h.transactionService.GroupByStatusTransaction()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": true, "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(transactions)
}

func (h *IncomingHandler) SelectStatusDetail(c *fiber.Ctx) error {
	transactions, err := h.transactionService.GroupByStatusDetailTransaction()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": true, "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(transactions)
}

func (h *IncomingHandler) SelectAdnet(c *fiber.Ctx) error {
	transactions, err := h.transactionService.GroupByAdnetTransaction()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": true, "message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(transactions)
}

func (h *IncomingHandler) ReportDaily(c *fiber.Ctx) error {
	transactions, err := h.transactionService.GroupByStatusTransaction()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": true, "message": err.Error()})
	}
	return c.Render("report/daily", fiber.Map{
		"transactions": transactions,
	})
}

func (h *IncomingHandler) AveragePerUser(c *fiber.Ctx) error {
	/**
	 * Body Parser
	 */
	req := new(entity.ReqArpuParams)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	subs, err := h.subscriptionService.AveragePerUser(req.GetStart(), req.GetEnd(), req.GetToRenew(), req.GetService())
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"code":  fiber.StatusOK,
		"data":  subs,
	})
}

func (h *IncomingHandler) Auth(c *fiber.Ctx) error {
	/**
	 * Body Parser
	 */
	req := new(entity.AuthPortalRequest)

	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusBadRequest,
				Message:    err.Error(),
			},
		)
	}

	errors := ValidateStruct(*req)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	if c.Get("Cf-Connecting-Ip") != "" {
		req.SetIpAddress(c.Get("Cf-Connecting-Ip"))
	} else {
		req.SetIpAddress(c.Get("X-Forwarded-For"))
	}

	if !h.serviceService.IsServiceByCategory(strings.ToUpper(c.Params("category"))) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "service_not_found",
			},
		)
	}

	if !h.subscriptionService.GetPinActiveSub(strings.ToUpper(c.Params("category")), req.GetPin()) {
		return c.Status(fiber.StatusNotFound).JSON(
			&model.WebResponse{
				Error:      true,
				StatusCode: fiber.StatusNotFound,
				Message:    "pin_not_found",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		&model.WebResponse{
			Error:      false,
			StatusCode: fiber.StatusOK,
			Message:    "authenticated",
			IpAddress:  req.GetIpAddress(),
		},
	)
}
