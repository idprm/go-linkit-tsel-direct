package entity

import (
	"database/sql"
	"net/url"
	"strings"
	"time"
)

const (
	MO_REG       = "REG"
	MO_UNREG     = "UNREG"
	MO_OFF       = "OFF"
	MO_CONFIRM   = "Y"
	VALID_PREFIX = "628"
)

type (
	ArrayReqSub struct {
		Req ReqSub `json:"request"`
	}

	ReqSub struct {
		Sms       string `json:"sms" form:"sms"`
		Msisdn    string `json:"msisdn" form:"msisdn"`
		Adn       string `json:"adn" form:"adn"`
		IpAddress string `json:"ip_address" form:"ip_address"`
	}
)

type ReqTrafficParams struct {
	TxId           string `json:"tx_id,omitempty"`
	ServiceId      int    `json:"service_id"`
	CampKeyword    string `json:"camp_keyword,omitempty"`
	CampSubKeyword string `json:"camp_sub_keyword,omitempty"`
	Adnet          string `json:"adnet,omitempty"`
	PubId          string `json:"pub_id,omitempty"`
	AffSub         string `json:"aff_sub,omitempty"`
	Browser        string `json:"browser,omitempty"`
	OS             string `json:"os,omitempty"`
	Device         string `json:"device,omitempty"`
	Referer        string `json:"referer,omitempty"`
	IpAddress      string `json:"ip_address,omitempty"`
}

func (e *ReqTrafficParams) GetTxId() string {
	return e.TxId
}

func (e *ReqTrafficParams) GetServiceId() int {
	return e.ServiceId
}

func (e *ReqTrafficParams) GetCampKeyword() string {
	return strings.ToUpper(e.CampKeyword)
}

func (e *ReqTrafficParams) GetCampSubKeyword() string {
	return strings.ToUpper(e.CampSubKeyword)
}

func (e *ReqTrafficParams) GetAdnet() string {
	return e.Adnet
}

func (e *ReqTrafficParams) GetPubId() string {
	return e.PubId
}

func (e *ReqTrafficParams) GetAffSub() string {
	return e.AffSub
}

func (e *ReqTrafficParams) GetBrowser() string {
	return e.Browser
}

func (e *ReqTrafficParams) GetOS() string {
	return e.OS
}

func (e *ReqTrafficParams) GetDevice() string {
	return e.Device
}

func (e *ReqTrafficParams) GetReferer() string {
	return e.Referer
}

func (e *ReqTrafficParams) GetIpAddress() string {
	return e.IpAddress
}

type ReqMOParams struct {
	SMS       string `validate:"required" query:"sms" json:"sms"`
	Adn       string `query:"adn" json:"adn"`
	Msisdn    string `validate:"required" query:"msisdn" json:"msisdn"`
	Channel   string `query:"channel" json:"channel"`
	TrxId     string `query:"trx_id" json:"trx_id"`
	Number    string `query:"http_segment_number" json:"http_segment_number"`
	Count     string `query:"http_segment_count" json:"http_segment_count"`
	IpAddress string `query:"ip_address" json:"ip_address"`
}

type ReqMOBody struct {
	MessageID struct {
		Sms struct {
			Retry struct {
				Count       string `json:"count" xml:"count"`
				Max         string `json:"max" xml:"max"`
				Destination struct {
					Address struct {
						Unknown struct {
							Cnpi string `json:"cnpi" xml:"cnpi"`
						} `json:"unknown" xml:"unknown"`
					} `json:"address" xml:"address"`
				} `json:"destination" xml:"destination"`
				Source struct {
					Address struct {
						Number struct {
							Type string `json:"type" xml:"type"`
						} `json:"number" xml:"number"`
					} `json:"address" xml:"address"`
				} `json:"source" xml:"source"`
				Ud struct {
					Type string `json:"type" xml:"type"`
				} `json:"ud" xml:"ud"`
				Param struct {
					Name  string `json:"name" xml:"name"`
					Value string `json:"value" xml:"value"`
				} `json:"param" xml:"param"`
			} `json:"retry"`
		} `json:"sms" xml:"sms"`
	} `json:"message" xml:"message"`
}

type ReqMTParams struct {
	SMS    string `url:"sms,omitempty" query:"sms"`
	CpId   string `url:"cpid,omitempty" query:"cpid"`
	Pwd    string `url:"pwd,omitempty" query:"pwd"`
	Msisdn string `url:"msisdn,omitempty" query:"msisdn"`
	TrxId  string `url:"trx_id,omitempty" query:"trx_id"`
	Sid    string `url:"sid,omitempty" query:"sid"`
	Sender string `url:"sender,omitempty" query:"sender"`
	Tid    string `url:"tid,omitempty" query:"tid"`
}

type ReqMTBody struct {
	Message struct {
		Sms struct {
			Type        string `xml:"type,attr"`
			Destination struct {
				Address struct {
					Number string `xml:"number"`
				} `xml:"address"`
			} `xml:"destination"`
			Source struct {
				Address struct {
					Number string `xml:"number"`
				} `xml:"address"`
			} `xml:"source"`
			Ud    string           `xml:"ud"`
			Param []ReqMTBodyParam `xml:"param"`
		} `xml:"sms"`
	} `xml:"message"`
}

type ReqNotifParams struct {
	Subscription *Subscription
	Service      *Service
	Action       string `json:"action"`
	Pin          string `json:"pin"`
}

type ReqPostbackParams struct {
	Verify       *Verify
	Subscription *Subscription
	Service      *Service
	Postback     *Postback
	Action       string `json:"action"`
	Status       string `json:"status"`
	AffSub       string `json:"aff_sub"`
	IsSuccess    bool   `json:"is_success"`
}

type ReqMTBodyParam struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type ReqOptInParam struct {
	Service        string `json:"service" query:"service"`
	Adnet          string `json:"adnet" query:"ad"`
	PubId          string `json:"pub_id" query:"pubid"`
	AffSub         string `json:"aff_sub" query:"aff_sub"`
	CampKeyword    string `json:"keyword" query:"keyword"`
	CampSubKeyword string `json:"subkey" query:"subkey"`
	IpAddress      string `json:"ip_address" query:"ip"`
}

func (r *ReqOptInParam) GetService() string {
	return strings.ToUpper(r.Service)
}

func (r *ReqOptInParam) GetAdnet() string {
	return r.Adnet
}

func (r *ReqOptInParam) GetPubId() string {
	return r.PubId
}

func (r *ReqOptInParam) GetAffSub() string {
	return r.AffSub
}

func (r *ReqOptInParam) GetCampKeyword() string {
	decodeUrl, _ := url.QueryUnescape(r.CampKeyword)
	return strings.ToUpper(decodeUrl)
}

func (r *ReqOptInParam) GetCampSubKeyword() string {
	return strings.ToUpper(r.CampSubKeyword)
}

func (r *ReqOptInParam) GetIpAddress() string {
	return r.IpAddress
}

func (r *ReqOptInParam) SetService(data string) {
	r.Service = data
}

func (r *ReqOptInParam) SetIpAddress(ip string) {
	r.IpAddress = ip
}

type SuccessQueryParamsRequest struct {
	Token string `query:"token"`
	TrxId string `query:"trx_id"`
}

func (e *SuccessQueryParamsRequest) GetToken() string {
	return e.Token
}

func (e *SuccessQueryParamsRequest) GetTrxId() string {
	return e.TrxId
}

type CampaignToolsRequest struct {
	Service        string `json:"srv" query:"srv"`
	Dynamic        string `json:"dyn" query:"dyn"`
	Adnet          string `json:"adnet" query:"ad"`
	PubId          string `json:"pub_id" query:"pubid"`
	AffSub         string `json:"aff_sub" query:"aff_sub"`
	CampKeyword    string `json:"keyword" query:"keyword"`
	CampSubKeyword string `json:"subkey" query:"subkey"`
	IpAddress      string `json:"ip_address" query:"ip"`
	OS             string `json:"os" query:"os"`
	Browser        string `json:"browser" query:"browser"`
	UA             string `json:"useragent" query:"useragent"`
	Referer        string `json:"referer" query:"referer"`
	GTag           string `json:"gtag" query:"gtag"`
	Tiktok         string `json:"tiktok" query:"tiktok"`
}

func (r *CampaignToolsRequest) GetService() string {
	message := strings.ToUpper(r.Service)
	index := strings.Split(message, " ")
	if len(index[0]) > 0 {
		return index[0]
	}
	return ""
}

func (r *CampaignToolsRequest) GetDynamic() string {
	message := strings.ToUpper(r.Dynamic)
	index := strings.Split(message, " ")
	if len(index[0]) > 0 {
		return index[0]
	}
	return ""
}

func (r *CampaignToolsRequest) GetSubKeyword() string {
	message := strings.ToUpper(r.Service)
	index := strings.Split(message, " ")
	if len(index) > 1 {
		return index[1]
	}
	return ""
}

func (r *CampaignToolsRequest) GetSubDynamic() string {
	message := strings.ToUpper(r.Dynamic)
	index := strings.Split(message, " ")
	if len(index) > 1 {
		return index[1]
	}
	return ""
}

func (r *CampaignToolsRequest) GetAdnet() string {
	return r.Adnet
}

func (r *CampaignToolsRequest) GetPubId() string {
	return r.PubId
}

func (r *CampaignToolsRequest) GetAffSub() string {
	return r.AffSub
}

func (r *CampaignToolsRequest) GetIpAddress() string {
	return r.IpAddress
}

func (r *CampaignToolsRequest) SetIpAddress(ip string) {
	r.IpAddress = ip
}

func (r *CampaignToolsRequest) GetReferer() string {
	return r.Referer
}

func (r *CampaignToolsRequest) SetReferer(value string) {
	r.Referer = value
}

func (r *CampaignToolsRequest) GetOS() string {
	return r.OS
}

func (r *CampaignToolsRequest) GetBrowser() string {
	return r.Browser
}

func (r *CampaignToolsRequest) GetUA() string {
	return r.UA
}

func (r *CampaignToolsRequest) GetGTag() string {
	return strings.ToUpper(r.GTag)
}

func (r *CampaignToolsRequest) GetTiktok() string {
	return strings.ToUpper(r.Tiktok)
}

func (r *CampaignToolsRequest) IsBillable() bool {
	return strings.Contains(r.GetSubKeyword(), "BLB") ||
		strings.Contains(r.GetSubKeyword(), "BIL") ||
		strings.Contains(r.GetSubKeyword(), "LNK")
}

func (r *CampaignToolsRequest) IsSam() bool {
	return r.GetSubKeyword() == "SAM"
}

func (r *CampaignToolsRequest) IsYlc() bool {
	return r.GetSubKeyword() == "YLC"
}

func (r *CampaignToolsRequest) IsBng() bool {
	return r.GetSubKeyword() == "BNG"
}

func (r *CampaignToolsRequest) IsRdr() bool {
	return r.GetSubKeyword() == "RDR"
}

func (r *CampaignToolsRequest) IsStar() bool {
	return r.GetSubKeyword() == "STAR"
}

func (r *CampaignToolsRequest) IsMxo() bool {
	return r.GetSubKeyword() == "MXO"
}

func (r *CampaignToolsRequest) IsStars() bool {
	return r.GetSubKeyword() == "STARS"
}

func (r *CampaignToolsRequest) IsUnt() bool {
	return r.GetSubKeyword() == "UNT"
}

type CampaignToolsResponse struct {
	StatusCode int    `json:"status_code" xml:"status_code"`
	Token      string `json:"token" xml:"token"`
	UrlPromo   string `json:"url_promo" xml:"url_promo"`
}

type DailypushBodyRequest struct {
	TxId           string `json:"tx_id,omitempty"`
	SubscriptionId int64  `json:"subscription_id,omitempty"`
	ServiceId      int    `json:"service_id"`
	Msisdn         string `json:"msisdn"`
	Channel        string `json:"channel,omitempty"`
	CampKeyword    string `json:"camp_keyword,omitempty"`
	CampSubKeyword string `json:"camp_sub_keyword,omitempty"`
	Adnet          string `json:"adnet,omitempty"`
	PubID          string `json:"pub_id,omitempty"`
	AffSub         string `json:"aff_sub,omitempty"`
	Subject        string `json:"subject,omitempty"`
	StatusCode     string `json:"status_code,omitempty"`
	StatusDetail   string `json:"status_detail,omitempty"`
	IsCharge       bool   `json:"is_charge"`
	IpAddress      string `json:"ip_address,omitempty"`
	Action         string `json:"action,omitempty"`
}

func (e *DailypushBodyRequest) SetAction(value string) {
	e.Action = value
}

func (e *DailypushBodyRequest) IsRenewal() bool {
	return e.Action == "RENEWAL"
}

func (e *DailypushBodyRequest) IsRetry() bool {
	return e.Action == "RETRY"
}

type ResponseMO struct {
	StatusCode int    `json:"status_code" xml:"status_code"`
	Message    string `json:"message" xml:"message"`
}

type ResponseDR struct {
	StatusCode int    `json:"status_code" xml:"status_code"`
	Message    string `json:"message" xml:"message"`
}

type ReqArpuParams struct {
	Start   string `json:"from" query:"from"`
	End     string `json:"to" query:"to"`
	ToRenew string `json:"to_renew" query:"renew"`
	Service string `json:"service" query:"service"`
}

func (e *ReqArpuParams) GetStart() string {
	return e.Start
}

func (e *ReqArpuParams) GetEnd() string {
	return e.End
}

func (e *ReqArpuParams) GetToRenew() string {
	return e.ToRenew
}

func (e *ReqArpuParams) GetService() string {
	return e.Service
}

type AveragePerUser struct {
	Name       string  `json:"name"`
	Service    string  `json:"service"`
	Adnet      string  `json:"adnet"`
	Subs       string  `json:"subs"`
	SubsActive string  `json:"subs_active"`
	Revenue    float64 `json:"revenue"`
}

type AveragePerUserResponse struct {
	Name       string `json:"name"`
	Service    string `json:"service"`
	Adnet      string `json:"adnet"`
	Subs       string `json:"subs"`
	SubsActive string `json:"subs_active"`
	Revenue    int    `json:"revenue"`
}

func (e *AveragePerUserResponse) SetRevenue(revenue float64) {
	e.Revenue = int(revenue)
}

type ErrorResponse struct {
	FailedField string `json:"failed_field" xml:"failed_field"`
	Tag         string `json:"tag" xml:"tag"`
	Value       string `json:"value" xml:"value"`
}

func NewReqMOParams(sms, adn, msisdn, channel string) *ReqMOParams {
	return &ReqMOParams{
		SMS:     sms,
		Adn:     adn,
		Msisdn:  msisdn,
		Channel: channel,
	}
}

func (s *ReqMOParams) GetSMS() string {
	return s.SMS
}

func (s *ReqMOParams) SetSMS(sms string) {
	s.SMS = strings.ToUpper(sms)
}

func (s *ReqMOParams) GetAdn() string {
	return s.Adn
}

func (s *ReqMOParams) GetMsisdn() string {
	return s.Msisdn
}

func (s *ReqMOParams) GetChannel() string {
	return s.Channel
}

func (s *ReqMOParams) GetIpAddress() string {
	return s.IpAddress
}

func (s *ReqMOParams) IsInValidPrefix() bool {
	return !strings.HasPrefix(s.Msisdn, VALID_PREFIX)
}

func (s *ReqMOParams) IsREG() bool {
	message := strings.ToUpper(s.SMS)
	index := strings.Split(message, " ")
	if index[0] == MO_REG && (strings.Contains(message, MO_REG)) {
		return true
	}
	return false
}

func (s *ReqMOParams) IsUNREG() bool {
	message := strings.ToUpper(s.SMS)
	index := strings.Split(message, " ")
	if index[0] == MO_UNREG && (strings.Contains(message, MO_UNREG)) {
		return true
	}
	if index[0] == MO_OFF && (strings.Contains(message, MO_OFF)) {
		return true
	}
	return false
}

func (s *ReqMOParams) IsOFF() bool {
	message := strings.ToUpper(s.SMS)
	index := strings.Split(message, " ")
	if index[0] == MO_OFF && (strings.Contains(message, MO_OFF)) {
		return true
	}
	return false
}

func (s *ReqMOParams) IsConfirm() bool {
	message := strings.ToUpper(s.SMS)
	index := strings.Split(message, " ")
	if index[0] == MO_CONFIRM && (strings.Contains(message, MO_CONFIRM)) {
		return true
	}
	return false
}

func (s *ReqMOParams) GetKeyword() string {
	return strings.ToUpper(s.SMS)
}

func (s *ReqMOParams) GetSubKeyword() string {
	message := strings.ToUpper(s.SMS)
	index := strings.Split(message, " ")

	if index[0] == MO_REG || index[0] == MO_UNREG || index[0] == MO_OFF {
		if strings.Contains(message, MO_REG) || strings.Contains(message, MO_UNREG) || strings.Contains(message, MO_OFF) {
			if len(index) > 1 {
				return index[1]
			}
			return ""
		}
		return ""
	}
	return ""
}

func (e *ReqNotifParams) IsSub() bool {
	return e.Action == "SUB"
}

func (e *ReqNotifParams) IsRenewal() bool {
	return e.Action == "RENEWAL"
}

func (e *ReqNotifParams) IsUnsub() bool {
	return e.Action == "UNSUB"
}

func (e *ReqPostbackParams) IsMO() bool {
	return e.Action == "MO"
}

func (e *ReqPostbackParams) IsMOUnsub() bool {
	return e.Action == "MO_UNSUB"
}

func (e *ReqPostbackParams) IsMT() bool {
	return e.Action == "MT"
}

// for retry firstpush
func (e *ReqPostbackParams) IsMTFirstpush() bool {
	return e.Action == "MT_FIRSTPUSH"
}

// for renewal dailypush & retry dailypush
func (e *ReqPostbackParams) IsMTDailypush() bool {
	return e.Action == "MT_DAILYPUSH"
}

func (e *ReqPostbackParams) GetIsSuccess() bool {
	return e.Status == "SUCCESS"
}

func (e *ReqPostbackParams) IsPostbackNotNull() bool {
	return e.Postback != nil
}

func (e *ReqPostbackParams) IsPostbackNull() bool {
	return e.Postback == nil
}

var formatDate = "2006-01-02T15:04:05Z07:00"

type SubscriptionToCSV struct {
	Country         string         `json:"country,omitempty"`
	Operator        string         `json:"operator,omitempty"`
	Service         string         `json:"service,omitempty"`
	Source          string         `json:"source,omitempty"`
	Msisdn          string         `json:"msisdn,omitempty"`
	LatestSubject   string         `json:"latest_subject,omitempty"`
	Cycle           string         `json:"cycle,omitempty"`
	Adnet           string         `json:"adnet,omitempty"`
	Revenue         string         `json:"revenue,omitempty"`
	SubsDate        sql.NullString `json:"subs_date,omitempty"`
	RenewalDate     sql.NullString `json:"renewal_date,omitempty"`
	FreemiumEndDate string         `json:"freemium_end_date,omitempty"`
	UnsubsFrom      string         `json:"unsubs_from,omitempty"`
	UnsubsDate      sql.NullString `json:"unsubs_date,omitempty"`
	ServicePrice    string         `json:"service_price,omitempty"`
	Currency        string         `json:"currency,omitempty"`
	ProfileStatus   string         `json:"profile_status,omitempty"`
	Publisher       string         `json:"publisher,omitempty"`
	Trxid           string         `json:"trxid,omitempty"`
	Pixel           string         `json:"pixel,omitempty"`
	Handset         string         `json:"handset,omitempty"`
	Browser         string         `json:"browser,omitempty"`
	AttemptCharging string         `json:"attempt_charging"`
	SuccessBilling  string         `json:"success_billing"`
	CampSubKeyword  string         `json:"camp_sub_keyword,omitempty"`
}

func (e *SubscriptionToCSV) SetLatestSubject(data string) {
	switch data {
	case "FIRSTPUSH":
		e.LatestSubject = "1"
	case "RENEWAL":
		e.LatestSubject = "0"
	case "UNSUB":
		e.LatestSubject = "-1"
	default:
		e.LatestSubject = "NA"
	}
}

func (e *SubscriptionToCSV) SetService(data, subkey string) {
	if subkey != "" {
		e.Service = data + " " + subkey
	} else {
		e.Service = data
	}
}

func (e *SubscriptionToCSV) SetSubsDate(data string) {
	dt, _ := time.Parse(formatDate, data)
	e.SubsDate.String = dt.Format("2006-01-02 15:04:05") + " +0700"
}

func (e *SubscriptionToCSV) SetRenewalDate(data string) {
	dt, _ := time.Parse(formatDate, data)
	e.RenewalDate.String = dt.Format("2006-01-02 15:04:05") + " +0700"
}

func (e *SubscriptionToCSV) SetUnsubsDate(data string) {
	dt, _ := time.Parse(formatDate, data)
	e.UnsubsDate.String = dt.Format("2006-01-02 15:04:05") + " +0700"
}

func (e *SubscriptionToCSV) SetProfileStatus(data string) {
	switch data {
	case "true":
		e.ProfileStatus = "active"
	case "false":
		e.ProfileStatus = "inactive"
	default:
		e.ProfileStatus = "NA"
	}
}

func (e *SubscriptionToCSV) SetAdnet(data string) {
	if data != "" {
		e.Adnet = data
	} else {
		e.Adnet = "NA"
	}
}

type TransactionToCSV struct {
	Country          string         `json:"country,omitempty"`
	Operator         string         `json:"operator,omitempty"`
	Service          string         `json:"service,omitempty"`
	Source           string         `json:"source,omitempty"`
	Msisdn           string         `json:"msisdn,omitempty"`
	Event            string         `json:"event,omitempty"`
	EventDate        sql.NullString `json:"even_date,omitempty"`
	Cycle            string         `json:"cycle,omitempty"`
	Revenue          string         `json:"revenue,omitempty"`
	ChargeDate       sql.NullString `json:"charge_date,omitempty"`
	Currency         string         `json:"currency,omitempty"`
	Publisher        string         `json:"publisher,omitempty"`
	Handset          string         `json:"handset,omitempty"`
	Browser          string         `json:"browser,omitempty"`
	TrxId            string         `json:"trxid,omitempty"`
	TelcoApiUrl      string         `json:"telco_api_url,omitempty"`
	TelcoApiResponse string         `json:"telco_api_response,omitempty"`
	SmsContent       string         `json:"sms_content,omitempty"`
	StatusSms        string         `json:"status_sms,omitempty"`
	CampSubKeyword   string         `json:"camp_sub_keyword,omitempty"`
}

func (e *TransactionToCSV) SetService(data, subkey string) {
	if subkey != "" {
		e.Service = data + " " + subkey
	} else {
		e.Service = data
	}
}

func (e *TransactionToCSV) SetEventDate(data string) {
	dt, _ := time.Parse(formatDate, data)
	e.EventDate.String = dt.Format("2006-01-02 15:04:05") + " +0700"
}
func (e *TransactionToCSV) SetChargeDate(data string) {
	dt, _ := time.Parse(formatDate, data)
	e.ChargeDate.String = dt.Format("2006-01-02 15:04:05") + " +0700"
}

func (e *TransactionToCSV) GetCycle() string {
	return e.Cycle
}

func (e *TransactionToCSV) SetCycle(data string) {
	switch data {
	case "7d":
		e.Cycle = "weekly"
	case "30d":
		e.Cycle = "monthly"
	default:
		e.Cycle = data
	}
}

type RabbitMQResponse struct {
	Messages int    `json:"messages"`
	Name     string `json:"name"`
}

func (r *RabbitMQResponse) IsRunning() bool {
	return r.Messages > 0
}

func (r *RabbitMQResponse) GetName() string {
	return r.Name
}

type AuthPortalRequest struct {
	PIN       string `validate:"required" json:"pin"`
	IpAddress string `json:"ip_address,omitempty"`
}

func (e *AuthPortalRequest) GetPin() string {
	return e.PIN
}

func (e *AuthPortalRequest) GetIpAddress() string {
	return e.IpAddress
}

func (m *AuthPortalRequest) SetIpAddress(data string) {
	m.IpAddress = data
}
