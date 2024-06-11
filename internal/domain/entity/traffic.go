package entity

import "time"

type TrafficCampaign struct {
	ID             int64     `json:"id,omitempty"`
	ServiceID      int       `json:"service_id,omitempty"`
	Service        *Service  `json:",omitempty"`
	CampKeyword    string    `json:"camp_keyword,omitempty"`
	CampSubKeyword string    `json:"camp_sub_keyword,omitempty"`
	Adnet          string    `json:"adnet,omitempty"`
	PubID          string    `json:"pub_id,omitempty"`
	AffSub         string    `json:"aff_sub,omitempty"`
	Browser        string    `json:"browser,omitempty"`
	OS             string    `json:"os,omitempty"`
	Device         string    `json:"device,omitempty"`
	IpAddress      string    `json:"ip_address,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

type TrafficMO struct {
	ID             int64     `json:"id,omitempty"`
	ServiceID      int       `json:"service_id,omitempty"`
	Service        *Service  `json:",omitempty"`
	Msisdn         string    `json:"msisdn"`
	Channel        string    `json:"channel,omitempty"`
	CampKeyword    string    `json:"camp_keyword,omitempty"`
	CampSubKeyword string    `json:"camp_sub_keyword,omitempty"`
	Adnet          string    `json:"adnet,omitempty"`
	PubID          string    `json:"pub_id,omitempty"`
	AffSub         string    `json:"aff_sub,omitempty"`
	IpAddress      string    `json:"ip_address,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}
