package entity

import (
	"net/url"
	"strings"
)

type Postback struct {
	ID         int    `json:"id"`
	SubKeyword string `json:"sub_keyword"`
	UrlMO      string `json:"url_mo"`
	UrlDN      string `json:"url_dn"`
	IsActive   bool   `json:"is_active"`
}

func (e *Postback) GetId() int {
	return e.ID
}

func (e *Postback) GetSubKeyword() string {
	return e.SubKeyword
}

func (e *Postback) GetUrlMO() string {
	return e.UrlMO
}

func (e *Postback) GetUrlDN() string {
	return e.UrlDN
}

func (e *Postback) SetUrlMO(sdc, msisdn, sms, trxid, trxdate string) {
	// mo?sdc={sdc}&msisdn={msisdn}&sms={sms}&trx_id={trxid}&trx_date={trxdate}
	replacer := strings.NewReplacer(
		"{sdc}", sdc,
		"{msisdn}", url.QueryEscape(msisdn),
		"{sms}", url.QueryEscape(sms),
		"{trxid}", trxid,
		"{trxdate}", trxdate)
	e.UrlMO = replacer.Replace(e.UrlMO)
}

func (e *Postback) SetUrlDN(sdc, msisdn, status, statusdesc, service, typeService, trxid, trxdate string) {
	// dn?sdc={sdc}&msisdn={msisdn}&status={status}&statusdesc={statusdesc}&service={service}&type={type}&trx_id={trxid}&trx_date={trxdate}
	replacer := strings.NewReplacer(
		"{sdc}", sdc,
		"{msisdn}", url.QueryEscape(msisdn),
		"{status}", url.QueryEscape(status),
		"{statusdesc}", url.QueryEscape(statusdesc),
		"{service}", url.QueryEscape(service),
		"{type}", url.QueryEscape(typeService),
		"{trxid}", trxid,
		"{trxdate}", trxdate)
	e.UrlDN = replacer.Replace(e.UrlDN)
}

func (e *Postback) IsSubKeyword(subkey string) bool {
	return e.SubKeyword == strings.ToUpper(subkey)
}
