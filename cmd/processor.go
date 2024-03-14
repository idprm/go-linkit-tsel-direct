package cmd

import (
	"database/sql"
	"encoding/json"
	"sync"

	"github.com/idprm/go-linkit-tsel/src/config"
	"github.com/idprm/go-linkit-tsel/src/domain/entity"
	"github.com/idprm/go-linkit-tsel/src/domain/repository"
	"github.com/idprm/go-linkit-tsel/src/handler"
	"github.com/idprm/go-linkit-tsel/src/logger"
	"github.com/idprm/go-linkit-tsel/src/services"
	"github.com/redis/go-redis/v9"
	"github.com/wiliehidayat87/rmqp"
)

type Processor struct {
	cfg    *config.Secret
	db     *sql.DB
	rmpq   rmqp.AMQP
	rdb    *redis.Client
	logger *logger.Logger
}

func NewProcessor(
	cfg *config.Secret,
	db *sql.DB,
	rmpq rmqp.AMQP,
	rdb *redis.Client,
	logger *logger.Logger,
) *Processor {
	return &Processor{
		cfg:    cfg,
		db:     db,
		rmpq:   rmpq,
		rdb:    rdb,
		logger: logger,
	}
}

func (p *Processor) MO(wg *sync.WaitGroup, message []byte) {
	/**
	 * -. Check Valid Prefix
	 * -. Filter REG / UNREG
	 * -. Check Blacklist
	 * -. Check Active Sub
	 * -. MT API
	 * -. Save Sub
	 * -/ Save Transaction
	 */
	campaignRepo := repository.NewCampaignRepository(p.db)
	campaignService := services.NewCampaignService(campaignRepo)
	blacklistRepo := repository.NewBlacklistRepository(p.db)
	blacklistService := services.NewBlacklistService(blacklistRepo)
	serviceRepo := repository.NewServiceRepository(p.db)
	serviceService := services.NewServiceService(serviceRepo)
	verifyRepo := repository.NewVerifyRepository(p.rdb)
	verifyService := services.NewVerifyService(verifyRepo)
	contentRepo := repository.NewContentRepository(p.db)
	contentService := services.NewContentService(contentRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(p.db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(p.db)
	transactionService := services.NewTransactionService(transactionRepo)
	historyRepo := repository.NewHistoryRepository(p.db)
	historyService := services.NewHistoryService(historyRepo)

	var req *entity.ReqMOParams
	json.Unmarshal([]byte(message), &req)
	reqMO := entity.NewReqMOParams(req.SMS, req.Adn, req.Msisdn, req.Channel)

	h := handler.NewMOHandler(
		p.cfg,
		p.rmpq,
		p.logger,
		campaignService,
		blacklistService,
		serviceService,
		verifyService,
		contentService,
		subscriptionService,
		transactionService,
		historyService,
		reqMO,
	)

	// check service by category
	if h.IsService() {
		// filter REG
		if reqMO.IsREG() {
			// filter not blacklist
			if !h.IsBlacklist() {
				// active sub
				if !h.IsActiveSub() {
					// Firstpush MT API
					h.Firstpush()
				} else {
					h.Logger(reqMO, "ALREADY_SUB")
				}
			} else {
				h.Logger(reqMO, "BLACKLIST")
			}
		}
		if reqMO.IsUNREG() {
			// active sub
			if h.IsActiveSub() {
				// unsub
				h.Unsub()
			} else {
				h.Logger(reqMO, "ALREADY_UNSUB")
			}
		}
		if reqMO.IsConfirm() {
			// confirm
			h.Confirm()
		}
	}

	wg.Done()
}

func (p *Processor) Renewal(wg *sync.WaitGroup, message []byte) {
	/**
	 * load repo
	 */
	serviceRepo := repository.NewServiceRepository(p.db)
	serviceService := services.NewServiceService(serviceRepo)
	contentRepo := repository.NewContentRepository(p.db)
	contentService := services.NewContentService(contentRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(p.db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(p.db)
	transactionService := services.NewTransactionService(transactionRepo)

	// parsing json to string
	var sub *entity.Subscription
	json.Unmarshal(message, &sub)

	h := handler.NewRenewalHandler(
		p.cfg,
		p.rmpq,
		p.logger,
		sub,
		serviceService,
		contentService,
		subscriptionService,
		transactionService,
	)

	// Dailypush MT API
	h.Dailypush()

	wg.Done()
}

func (p *Processor) RetryFp(wg *sync.WaitGroup, message []byte) {
	/**
	 * load repo
	 */
	serviceRepo := repository.NewServiceRepository(p.db)
	serviceService := services.NewServiceService(serviceRepo)
	contentRepo := repository.NewContentRepository(p.db)
	contentService := services.NewContentService(contentRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(p.db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(p.db)
	transactionService := services.NewTransactionService(transactionRepo)

	// parsing json to string
	var sub *entity.Subscription
	json.Unmarshal(message, &sub)

	h := handler.NewRetryHandler(p.cfg, p.rmpq, p.logger, sub, serviceService, contentService, subscriptionService, transactionService)

	h.Firstpush()

	wg.Done()
}

func (p *Processor) RetryDp(wg *sync.WaitGroup, message []byte) {
	/**
	 * load repo
	 */
	serviceRepo := repository.NewServiceRepository(p.db)
	serviceService := services.NewServiceService(serviceRepo)
	contentRepo := repository.NewContentRepository(p.db)
	contentService := services.NewContentService(contentRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(p.db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(p.db)
	transactionService := services.NewTransactionService(transactionRepo)

	// parsing json to string
	var sub *entity.Subscription
	json.Unmarshal(message, &sub)

	h := handler.NewRetryHandler(p.cfg, p.rmpq, p.logger, sub, serviceService, contentService, subscriptionService, transactionService)

	h.Dailypush()

	wg.Done()
}

func (p *Processor) RetryInsuff(wg *sync.WaitGroup, message []byte) {
	/**
	 * load repo
	 */
	serviceRepo := repository.NewServiceRepository(p.db)
	serviceService := services.NewServiceService(serviceRepo)
	contentRepo := repository.NewContentRepository(p.db)
	contentService := services.NewContentService(contentRepo)
	subscriptionRepo := repository.NewSubscriptionRepository(p.db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	transactionRepo := repository.NewTransactionRepository(p.db)
	transactionService := services.NewTransactionService(transactionRepo)

	// parsing json to string
	var sub *entity.Subscription
	json.Unmarshal(message, &sub)

	h := handler.NewRetryHandler(p.cfg, p.rmpq, p.logger, sub, serviceService, contentService, subscriptionService, transactionService)

	if sub.IsFirstpush() {
		if sub.IsRetryAtToday() {
			h.Firstpush()
		} else {
			h.Dailypush()
		}
	} else {
		h.Dailypush()
	}

	wg.Done()
}

func (p *Processor) Notif(wg *sync.WaitGroup, message []byte) {
	var req *entity.ReqNotifParams
	json.Unmarshal(message, &req)

	h := handler.NewNotifHandler(p.cfg, p.logger, req)

	if req.IsSub() {
		h.Sub()
	}

	if req.IsRenewal() {
		h.Renewal()
	}

	if req.IsUnsub() {
		h.Unsub()
	}

	wg.Done()
}

func (p *Processor) PostbackMO(wg *sync.WaitGroup, message []byte) {
	var req *entity.ReqPostbackParams
	json.Unmarshal(message, &req)

	h := handler.NewPostbackHandler(p.cfg, p.logger, req)

	if req.IsMO() {
		if req.Verify.IsSam() {
			h.SamMO()
		}
		if req.Verify.IsYlc() {
			h.YlcMO(req.Verify.GetAffSub())
		}

		if req.Verify.IsFs() {
			h.FsMO()
		}
		// non billable
		if !req.Verify.GetIsBillable() {
			if !req.Verify.IsSam() && !req.Verify.IsYlc() && !req.Verify.IsBng() && !req.Verify.IsFs() && !req.Verify.IsRdr() && !req.Verify.IsV2Test() {
				h.Postback()
			}
		}

		if req.Verify.IsV2Test() {
			h.PbV2Test()
		}
	}

	if req.IsMOUnsub() {
		if req.Subscription.IsSAM() {
			h.SamMOUnsub()
		}
	}

	if req.IsMT() {
		if req.IsSuccess {
			if req.Verify.GetIsBillable() {
				// if success charge hit pb billable
				h.Billable()
			}
			if req.Verify.IsYlc() {
				// if success charge hit pb ylc
				h.YlcMT(req.Verify.GetAffSub())
			}
		}

		if req.Verify.IsSam() {
			h.SamDN(req.Status)
		}

		if req.Verify.IsFs() {
			h.FsDN(req.Status)
		}
	}

	wg.Done()
}

func (p *Processor) PostbackMT(wg *sync.WaitGroup, message []byte) {
	var req *entity.ReqPostbackParams
	json.Unmarshal(message, &req)

	h := handler.NewPostbackHandler(p.cfg, p.logger, req)

	/**
	 * Renewal & Retry Dailypush
	 */
	if req.IsMTDailypush() {
		if req.Subscription.IsSAM() {
			h.SamDN(req.Status)
		}
		if req.Subscription.IsFs() {
			h.FsDN(req.Status)
		}

		if req.Subscription.IsYLC() {
			h.YlcMT(req.AffSub)
		}
	}

	/**
	 * Retry Firstpush
	 */
	if req.IsMTFirstpush() {
		if req.Subscription.IsSAM() {
			h.SamDN(req.Status)
		}
		if req.Subscription.IsYLC() {
			h.YlcMT(req.AffSub)
		}
		if req.Subscription.IsFs() {
			h.FsDN(req.Status)
		}
	}

	wg.Done()
}
