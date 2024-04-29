package entity

import (
	"strings"
	"time"
)

type Subscription struct {
	ID                   int64     `json:"id"`
	ServiceID            int       `json:"service_id,omitempty"`
	Service              *Service  `json:"service,omitempty"`
	Category             string    `json:"category,omitempty"`
	Msisdn               string    `json:"msisdn"`
	Channel              string    `json:"channel,omitempty"`
	CampKeyword          string    `json:"camp_keyword,omitempty"`
	CampSubKeyword       string    `json:"camp_sub_keyword,omitempty"`
	Adnet                string    `json:"adnet,omitempty"`
	PubID                string    `json:"pub_id,omitempty"`
	AffSub               string    `json:"aff_sub,omitempty"`
	LatestTrxId          string    `json:"trx_id,omitempty"`
	LatestKeyword        string    `json:"latest_keyword,omitempty"`
	LatestSubject        string    `json:"latest_subject,omitempty"`
	LatestStatus         string    `json:"latest_status,omitempty"`
	LatestPIN            string    `json:"latest_pin,omitempty"`
	LatestPayload        string    `json:"latest_payload,omitempty"`
	Amount               float64   `json:"amount,omitempty"`
	TrialAt              time.Time `json:"trial_at,omitempty"`
	RenewalAt            time.Time `json:"renewal_at,omitempty"`
	UnsubAt              time.Time `json:"unsub_at,omitempty"`
	ChargeAt             time.Time `json:"charge_at,omitempty"`
	RetryAt              time.Time `json:"retry_at,omitempty"`
	PurgeAt              time.Time `json:"purge_at,omitempty"`
	PurgeReason          string    `json:"purge_reason,omitempty"`
	Success              uint      `json:"success,omitempty"`
	Failed               uint      `json:"failed,omitempty"`
	IpAddress            string    `json:"ip_address,omitempty"`
	TotalFirstpush       uint      `json:"total_firstpush,omitempty"`
	TotalRenewal         uint      `json:"total_renewal,omitempty"`
	TotalAmountFirstpush float64   `json:"total_amount_firstpush,omitempty"`
	TotalAmountRenewal   float64   `json:"total_amount_renewal,omitempty"`
	ChargingCount        uint      `json:"charging_count,omitempty"`
	ChargingCountAll     uint      `json:"charging_count_all,omitempty"`
	IsTrial              bool      `json:"is_trial,omitempty"`
	IsRetry              bool      `json:"is_retry"`
	IsConfirm            bool      `json:"is_confirm"`
	IsPurge              bool      `json:"is_purge"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (s *Subscription) GetServiceId() int {
	return s.ServiceID
}

func (s *Subscription) GetCategory() string {
	return s.Category
}

func (s *Subscription) GetMsisdn() string {
	return s.Msisdn
}

func (s *Subscription) GetChannel() string {
	return s.Channel
}

func (s *Subscription) GetCampKeyword() string {
	return s.CampKeyword
}

func (s *Subscription) GetCampSubKeyword() string {
	return strings.ToUpper(s.CampSubKeyword)
}

func (s *Subscription) IsSAM() bool {
	return s.GetCampSubKeyword() == "SAM"
}

func (s *Subscription) IsYLC() bool {
	return s.GetCampSubKeyword() == "YLC" || s.GetCampSubKeyword() == "YL2"
}

func (s *Subscription) IsBng() bool {
	return s.GetCampSubKeyword() == "BNG"
}

func (s *Subscription) IsFs() bool {
	return s.GetCampSubKeyword() == "FS"
}

func (s *Subscription) IsRdr() bool {
	return s.GetCampSubKeyword() == "RDR"
}

func (s *Subscription) IsPlw() bool {
	return s.GetCampSubKeyword() == "PLW"
}

func (s *Subscription) IsStar() bool {
	return s.GetCampSubKeyword() == "STAR"
}

func (s *Subscription) GetAdnet() string {
	return s.Adnet
}

func (s *Subscription) GetPubId() string {
	return s.PubID
}

func (s *Subscription) GetAffSub() string {
	return s.AffSub
}

func (s *Subscription) GetLatestTrxId() string {
	return s.LatestTrxId
}

func (s *Subscription) GetLatestKeyword() string {
	return s.LatestKeyword
}

func (s *Subscription) GetLatestSubject() string {
	return strings.ToUpper(s.LatestSubject)
}

func (s *Subscription) GetLatestStatus() string {
	return s.LatestStatus
}

func (s *Subscription) GetLatestPIN() string {
	return s.LatestPIN
}

func (s *Subscription) IsLatestPIN() bool {
	return s.LatestPIN != ""
}

func (s *Subscription) GetIpAddress() string {
	return s.IpAddress
}

func (s *Subscription) GetChargingCount() uint {
	return s.ChargingCount
}

func (s *Subscription) GetChargingcountAll() uint {
	return s.ChargingCountAll
}

func (s *Subscription) SetIsActive(active bool) {
	s.IsActive = active
}

func (s *Subscription) SetIsConfirm(confirm bool) {
	s.IsConfirm = confirm
}

func (s *Subscription) SetIsRetry(retry bool) {
	s.IsRetry = retry
}

func (s *Subscription) SetIsTrial(trial bool) {
	s.IsTrial = trial
}

func (s *Subscription) SetRenewalAt(renewalAt time.Time) {
	s.RenewalAt = renewalAt
}

func (s *Subscription) SetRetryAt(retryAt time.Time) {
	s.RetryAt = retryAt
}

func (s *Subscription) SetChargeAt(chargeAt time.Time) {
	s.ChargeAt = chargeAt
}

func (s *Subscription) SetUnsubAt(unsubAt time.Time) {
	s.UnsubAt = unsubAt
}

func (s *Subscription) SetLatestSubject(latestSubject string) {
	s.LatestSubject = latestSubject
}

func (s *Subscription) SetLatestStatus(latestStatus string) {
	s.LatestStatus = latestStatus
}

func (s *Subscription) SetChannel(channel string) {
	s.Channel = channel
}

func (s *Subscription) SetAdnet(adnet string) {
	s.Adnet = adnet
}

func (s *Subscription) SetPubID(pubId string) {
	s.PubID = pubId
}

func (s *Subscription) SetAffSub(affsub string) {
	s.AffSub = affsub
}

func (s *Subscription) IsCreatedAtToday() bool {
	return s.CreatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02")
}

func (s *Subscription) IsRetryAtToday() bool {
	return s.RetryAt.Format("2006-01-02") == time.Now().Format("2006-01-02")
}

func (s *Subscription) IsFirstpush() bool {
	return s.GetLatestSubject() == "FIRSTPUSH"
}

func (s *Subscription) IsRenewal() bool {
	return s.GetLatestSubject() == "RENEWAL"
}
