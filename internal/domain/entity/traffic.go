package entity

import (
	"strings"
	"time"
)

type TrafficCampaign struct {
	ID             int64     `json:"id,omitempty"`
	TxId           string    `json:"tx_id,omitempty"`
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
	Referer        string    `json:"referer,omitempty"`
	IpAddress      string    `json:"ip_address,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

func (e *TrafficCampaign) GetServiceId() int {
	return e.ServiceID
}

func (e *TrafficCampaign) GetTxId() string {
	return e.TxId
}

func (e *TrafficCampaign) GetCampKeyword() string {
	return strings.ToUpper(e.CampKeyword)
}

func (e *TrafficCampaign) GetCampSubKeyword() string {
	return strings.ToUpper(e.CampSubKeyword)
}

func (e *TrafficCampaign) GetAdnet() string {
	return strings.ToUpper(e.Adnet)
}

func (e *TrafficCampaign) GetPubID() string {
	return strings.ToUpper(e.PubID)
}

func (e *TrafficCampaign) GetAffSub() string {
	return strings.ToUpper(e.AffSub)
}

func (e *TrafficCampaign) GetBrowser() string {
	return e.Browser
}

func (e *TrafficCampaign) GetOS() string {
	return e.OS
}

func (e *TrafficCampaign) GetDevice() string {
	return e.Device
}

func (e *TrafficCampaign) GetReferer() string {
	return e.Referer
}

func (e *TrafficCampaign) GetIpAddress() string {
	return e.IpAddress
}

func (e *TrafficCampaign) SetTxId(value string) {
	e.TxId = value
}

func (e *TrafficCampaign) SetReferer(value string) {
	e.Referer = value
}

type TrafficMO struct {
	ID             int64     `json:"id,omitempty"`
	TxId           string    `json:"tx_id,omitempty"`
	ServiceID      int       `json:"service_id,omitempty"`
	Service        *Service  `json:",omitempty"`
	Msisdn         string    `json:"msisdn"`
	Channel        string    `json:"channel,omitempty"`
	CampKeyword    string    `json:"camp_keyword,omitempty"`
	CampSubKeyword string    `json:"camp_sub_keyword,omitempty"`
	Subject        string    `json:"subject,omitempty"`
	Adnet          string    `json:"adnet,omitempty"`
	PubID          string    `json:"pub_id,omitempty"`
	AffSub         string    `json:"aff_sub,omitempty"`
	IpAddress      string    `json:"ip_address,omitempty"`
	IsCharge       bool      `json:"is_charge"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

func (e *TrafficMO) GetTxId() string {
	return e.TxId
}

func (e *TrafficMO) GetServiceId() int {
	return e.ServiceID
}

func (e *TrafficMO) GetMsisdn() string {
	return e.Msisdn
}

func (e *TrafficMO) GetCampKeyword() string {
	return strings.ToUpper(e.CampKeyword)
}

func (e *TrafficMO) GetCampSubKeyword() string {
	return strings.ToUpper(e.CampSubKeyword)
}

func (e *TrafficMO) GetSubject() string {
	return e.Subject
}

func (e *TrafficMO) GetAdnet() string {
	return strings.ToUpper(e.Adnet)
}

func (e *TrafficMO) GetPubId() string {
	return strings.ToUpper(e.PubID)
}

func (e *TrafficMO) GetAffSub() string {
	return strings.ToUpper(e.AffSub)
}

func (e *TrafficMO) GetIpAddress() string {
	return e.IpAddress
}

func (e *TrafficMO) SetTxId(value string) {
	e.TxId = value
}
