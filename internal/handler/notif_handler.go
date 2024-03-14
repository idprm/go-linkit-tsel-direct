package handler

import (
	"encoding/json"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/idprm/go-linkit-tsel/internal/providers/portal"
)

type NotifHandler struct {
	logger *logger.Logger
	req    *entity.ReqNotifParams
}

func NewNotifHandler(
	logger *logger.Logger,
	req *entity.ReqNotifParams,
) *NotifHandler {

	return &NotifHandler{
		logger: logger,
		req:    req,
	}
}

func (h *NotifHandler) Sub() {
	p := portal.NewPortal(h.logger, h.req.Subscription, h.req.Service, h.req.Pin, "success")
	p.Subscription()
}

func (h *NotifHandler) Renewal() {
	p := portal.NewPortal(h.logger, h.req.Subscription, h.req.Service, h.req.Pin, "success")
	notif, err := p.Renewal()
	if err != nil {
		h.logger.Writer(err)
	}
	/**
	 *  Parsing Response Notif Renewal
	 */
	type resRenewal struct {
		Success int    `json:"success"`
		Message string `json:"message"`
	}
	var responseRenewal resRenewal
	json.Unmarshal(notif, &responseRenewal)

	if responseRenewal.Message != "successfully renewal" || responseRenewal.Success != 1 {
		p.Subscription()
	}
}

func (h *NotifHandler) Unsub() {
	p := portal.NewPortal(h.logger, h.req.Subscription, h.req.Service, "", "")
	p.Unsubscription()
}
