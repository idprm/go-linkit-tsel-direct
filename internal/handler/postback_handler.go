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
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Verify.IsCampTool)
	p.Send()
}

func (h *PostbackHandler) Billable() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Verify.IsCampTool)
	p.Billable()
}

func (h *PostbackHandler) SamMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.SamMO()
}

func (h *PostbackHandler) SamMOUnsub() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.SamMOUnsub()
}

func (h *PostbackHandler) YlcMO(affSub string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.YlcMO(affSub)
}

func (h *PostbackHandler) FsMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.FsMO()
}

func (h *PostbackHandler) SamDN(status string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.SamDN(status)
}

func (h *PostbackHandler) YlcMT(affSub string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.YlcMT(affSub)
}

func (h *PostbackHandler) FsDN(status string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.FsDN(status)
}

func (h *PostbackHandler) PbV2Test() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, h.req.Verify.IsCampTool)
	p.SendTestV2()
}

func (h *PostbackHandler) PlwMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.PlwMO()
}

func (h *PostbackHandler) PlwMOUnsub() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.PlwMOUnsub()
}

func (h *PostbackHandler) PlwDN(status string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.PlwDN(status)
}

func (h *PostbackHandler) StarMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.StarMO()
}

func (h *PostbackHandler) MxoMO() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.MxoMO()
}

func (h *PostbackHandler) MxoMOUnsub() {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.MxoMOUnsub()
}

func (h *PostbackHandler) MxoDN(status string) {
	p := postback.NewPostback(h.logger, h.req.Subscription, h.req.Service, false)
	p.MxoDN(status)
}
