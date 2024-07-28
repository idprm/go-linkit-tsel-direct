package entity

import "time"

type Dailypush struct {
	ID             int64         `json:"id,omitempty"`
	TxId           string        `json:"tx_id,omitempty"`
	SubscriptionID int64         `json:"subscription_id,omitempty"`
	Subscription   *Subscription `json:",omitempty"`
	ServiceID      int           `json:"service_id,omitempty"`
	Service        *Service      `json:",omitempty"`
	Msisdn         string        `json:"msisdn"`
	Channel        string        `json:"channel,omitempty"`
	CampKeyword    string        `json:"camp_keyword,omitempty"`
	CampSubKeyword string        `json:"camp_sub_keyword,omitempty"`
	Adnet          string        `json:"adnet,omitempty"`
	PubID          string        `json:"pub_id,omitempty"`
	AffSub         string        `json:"aff_sub,omitempty"`
	Subject        string        `json:"subject,omitempty"`
	StatusCode     string        `json:"status_code,omitempty"`
	StatusDetail   string        `json:"status_detail,omitempty"`
	IsCharge       bool          `json:"is_charge"`
	IpAddress      string        `json:"ip_address,omitempty"`
	CreatedAt      time.Time     `json:"created_at,omitempty"`
	UpdatedAt      time.Time     `json:"updated_at,omitempty"`
}
