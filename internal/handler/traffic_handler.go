package handler

import (
	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/services"
)

type TrafficHandler struct {
	trafficService services.ITrafficService
	req            *entity.ReqTrafficParams
}

func NewTrafficHandler(
	trafficService services.ITrafficService,
	req *entity.ReqTrafficParams,
) *TrafficHandler {
	return &TrafficHandler{
		trafficService: trafficService,
		req:            req,
	}
}

func (h *TrafficHandler) Campaign() {
	h.trafficService.SaveCampaign(
		&entity.TrafficCampaign{
			ServiceID:      h.req.ServiceId,
			CampKeyword:    h.req.CampKeyword,
			CampSubKeyword: h.req.CampSubKeyword,
			Adnet:          h.req.Adnet,
			PubID:          h.req.PubId,
			AffSub:         h.req.AffSub,
			Browser:        h.req.Browser,
			OS:             h.req.OS,
			Device:         h.req.Device,
			IpAddress:      h.req.IpAddress,
		},
	)
}
