package handler

import (
	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/services"
)

type DailypushHandler struct {
	dailypushService services.IDailypushService
	req              *entity.DailypushBodyRequest
}

func NewDailypushHandler(
	dailypushService services.IDailypushService,
	req *entity.DailypushBodyRequest,
) *DailypushHandler {
	return &DailypushHandler{
		dailypushService: dailypushService,
		req:              req,
	}
}

func (h *DailypushHandler) Dailypush() {
	if h.req.IsRenewal() {
		h.dailypushService.Save(
			&entity.Dailypush{
				TxId:           h.req.TxId,
				SubscriptionID: h.req.SubscriptionId,
				ServiceID:      h.req.ServiceId,
				Msisdn:         h.req.Msisdn,
				Channel:        h.req.Channel,
				CampKeyword:    h.req.CampKeyword,
				CampSubKeyword: h.req.CampSubKeyword,
				Adnet:          h.req.Adnet,
				PubID:          h.req.PubID,
				AffSub:         h.req.AffSub,
				Subject:        h.req.Subject,
				StatusCode:     h.req.StatusCode,
				StatusDetail:   h.req.StatusDetail,
				IsCharge:       h.req.IsCharge,
				IpAddress:      h.req.IpAddress,
			},
		)
	}

	if h.req.IsRetry() {
		h.dailypushService.Update(
			&entity.Dailypush{
				TxId:           h.req.TxId,
				SubscriptionID: h.req.SubscriptionId,
				ServiceID:      h.req.ServiceId,
				Msisdn:         h.req.Msisdn,
				Channel:        h.req.Channel,
				CampKeyword:    h.req.CampKeyword,
				CampSubKeyword: h.req.CampSubKeyword,
				Adnet:          h.req.Adnet,
				PubID:          h.req.PubID,
				AffSub:         h.req.AffSub,
				Subject:        h.req.Subject,
				StatusCode:     h.req.StatusCode,
				StatusDetail:   h.req.StatusDetail,
				IsCharge:       h.req.IsCharge,
				IpAddress:      h.req.IpAddress,
			},
		)
	}
}
