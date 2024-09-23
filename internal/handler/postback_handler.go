package handler

import (
	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/idprm/go-linkit-tsel/internal/providers/postback"
)

type PostbackHandler struct {
	logger *logger.Logger
	req    *entity.ReqPostbackParams
}

func NewPostbackHandler(
	logger *logger.Logger,
	req *entity.ReqPostbackParams,
) *PostbackHandler {
	return &PostbackHandler{
		logger: logger,
		req:    req,
	}
}

func (h *PostbackHandler) Postback() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, h.req.Verify.IsCampTool)
	p.Send()
}

func (h *PostbackHandler) Billable() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, h.req.Verify.IsCampTool)
	p.Billable()
}

func (h *PostbackHandler) SamMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.SamMO()
}

func (h *PostbackHandler) SamMOUnsub() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.SamMOUnsub()
}

func (h *PostbackHandler) YlcMO(affSub string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.YlcMO(affSub)
}

func (h *PostbackHandler) FsMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.FsMO()
}

func (h *PostbackHandler) SamDN(status string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.SamDN(status)
}

func (h *PostbackHandler) YlcMT(affSub string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.YlcMT(affSub)
}

func (h *PostbackHandler) FsDN(status string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.FsDN(status)
}

func (h *PostbackHandler) PbV2Test() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, h.req.Verify.IsCampTool)
	p.SendTestV2()
}

func (h *PostbackHandler) PlwMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.PlwMO()
}

func (h *PostbackHandler) PlwMOUnsub() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.PlwMOUnsub()
}

func (h *PostbackHandler) PlwDN(status string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.PlwDN()
	if h.req.Subscription.IsSuccess() {
		p.PlwNotif(status)
	}
}

func (h *PostbackHandler) StarMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.StarMO()
}

func (h *PostbackHandler) MxoMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.MxoMO()
}

func (h *PostbackHandler) MxoMOUnsub() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.MxoMOUnsub()
}

func (h *PostbackHandler) MxoDN(status string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.MxoDN(status)
}

func (h *PostbackHandler) StarsMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.StarsMO()
}

func (h *PostbackHandler) UntMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.UntMO()
}

func (h *PostbackHandler) UntMOUnsub() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.UntMOUnsub()
}

func (h *PostbackHandler) UntDN(status string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.UntDN()
}

func (h *PostbackHandler) ExternalTrackerMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.MO()
}

func (h *PostbackHandler) ExternalTrackerMOUnsub() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.MOUnsub()
}

func (h *PostbackHandler) ExternalTrackerDN() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.DN()
}

func (h *PostbackHandler) PostbackFP() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Postback, false)
	p.FP()
}
