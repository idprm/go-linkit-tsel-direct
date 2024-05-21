package postback

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/idprm/go-linkit-tsel/internal/utils/response_utils"
	"github.com/idprm/go-linkit-tsel/internal/utils/uuid_utils"
	"github.com/sirupsen/logrus"
)

type Postback struct {
	logger       *logger.Logger
	subscription *entity.Subscription
	service      *entity.Service
	isCampTool   bool
}

func NewPostback(
	logger *logger.Logger,
	subscription *entity.Subscription,
	service *entity.Service,
	isCampTool bool,
) *Postback {
	return &Postback{
		logger:       logger,
		subscription: subscription,
		service:      service,
		isCampTool:   isCampTool,
	}
}

func (p *Postback) Send() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("partner", "linkittisel")
	q.Add("px", p.subscription.Adnet)
	if p.isCampTool {
		q.Add("serv_id", p.service.GetCode()+" "+p.subscription.GetCampSubKeyword())
		q.Add("token", response_utils.ParseToken(p.subscription.GetLatestKeyword()))
	} else {
		q.Add("serv_id", p.service.GetCode())
	}
	q.Add("msisdn", p.subscription.GetMsisdn())
	q.Add("trxid", p.subscription.GetLatestTrxId())
	q.Add("time", time.Now().String())

	req, err := http.NewRequest("GET", p.service.UrlPostback+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostback + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.Msisdn,
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK")

	return body, nil
}

func (p *Postback) SendTestV2() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	if p.isCampTool {
		q.Add("token", response_utils.ParseToken(p.subscription.GetLatestKeyword()))
	} else {
		q.Add("serv_id", p.service.GetCode())
	}
	q.Add("msisdn", p.subscription.GetMsisdn())
	q.Add("trxid", p.subscription.GetLatestTrxId())
	q.Add("time", time.Now().String())

	req, err := http.NewRequest("GET", p.service.UrlPostbackFsMO+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostbackFsMO + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK")

	return body, nil
}

func (p *Postback) Billable() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("partner", "linkittiselbil")
	q.Add("px", p.subscription.GetAdnet())
	if p.isCampTool {
		q.Add("serv_id", p.service.GetCode()+" "+p.subscription.GetCampSubKeyword())
		q.Add("token", response_utils.ParseToken(p.subscription.GetLatestKeyword()))
	} else {
		q.Add("serv_id", p.service.GetCode())
	}
	q.Add("msisdn", p.subscription.GetMsisdn())
	q.Add("trxid", p.subscription.GetLatestTrxId())
	q.Add("time", time.Now().String())

	req, err := http.NewRequest("GET", p.service.UrlPostbackBillable+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostbackBillable + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("BILLABLE")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.Msisdn,
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("BILLABLE")

	return body, nil
}

/**
 * Message Originated (SAM)
 */
func (p *Postback) SamMO() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("msisdn", p.subscription.GetMsisdn())

	if p.service.IsCloudplay() {
		q.Add("operator", "183")
		q.Add("id_service", "2131")
	}

	if p.service.IsGalays() {
		q.Add("operator", "198")
		q.Add("id_service", "2252")
	}
	// msisdn, id_service, operator, sms, trx_id, service_type, sdc, trx_date

	q.Add("sdc", "97770")
	q.Add("sms", p.subscription.GetCampKeyword()+" "+p.subscription.GetCampSubKeyword()+" "+p.subscription.GetAffSub())
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	q.Add("service_type", "2")
	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.UrlPostbackSamMO+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostbackSamMO + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_SAM_MO")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_SAM_MO")

	return body, nil
}

func (p *Postback) SamMOUnsub() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("msisdn", p.subscription.GetMsisdn())

	if p.service.IsCloudplay() {
		q.Add("operator", "183")
		q.Add("id_service", "2131")
	}

	if p.service.IsGalays() {
		q.Add("operator", "198")
		q.Add("id_service", "2252")
	}

	// msisdn, id_service, operator, sms, trx_id, service_type, sdc, trx_date

	q.Add("sdc", "97770")
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	q.Add("sms", p.subscription.GetLatestKeyword()+" "+p.subscription.GetCampSubKeyword())
	q.Add("service_type", "2")

	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.UrlPostbackSamMO+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostbackSamMO + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_SAM_MO_UNSUB")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_SAM_MO_UNSUB")

	return body, nil
}

/**
 * Delivery Notification (SAM)
 */
func (p *Postback) SamDN(status string) ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("msisdn", p.subscription.GetMsisdn())

	if p.service.IsCloudplay() {
		q.Add("operator", "183")
		q.Add("id_service", "2131")
	}

	if p.service.IsGalays() {
		q.Add("operator", "198")
		q.Add("id_service", "2252")
	}

	if status != "SUCCESS" {
		q.Add("status", "0")
	} else {
		q.Add("status", "1")
	}

	// msisdn, id_service, operator, trx_id, status, statusdesc, sdc, trx_date

	q.Add("statusdesc", strings.ToLower(status))
	q.Add("sdc", "97770")
	q.Add("service", p.subscription.GetCampKeyword()+" "+p.subscription.GetCampSubKeyword())
	q.Add("type", strings.ToLower(p.subscription.GetLatestSubject()))
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.UrlPostbackSamDN+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "request": p.service.UrlPostbackSamDN + "?" + q.Encode(), "trx_id": trxId}).Info("POSTBACK_SAM_DN")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_SAM_DN")

	return body, nil
}

func (p *Postback) YlcMO(affSub string) ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("transaction_id", affSub)

	req, err := http.NewRequest("GET", p.service.UrlPostbackYlcMO+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostbackYlcMO + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_YLC_MO")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_YLC_MO")

	return body, nil
}

func (p *Postback) YlcMT(affSub string) ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("type", "mt")
	q.Add("transaction_id", affSub)

	req, err := http.NewRequest("GET", p.service.UrlPostbackYlcMT+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostbackYlcMT + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_YLC_MT")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.Msisdn,
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_YLC_MT")

	return body, nil
}

/**
 * Message Originated (FS)
 */
func (p *Postback) FsMO() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("msisdn", p.subscription.GetMsisdn())
	q.Add("operator", "321")
	q.Add("sdc", "97770")
	q.Add("sms", p.subscription.GetCampKeyword()+" "+p.subscription.GetCampSubKeyword()+" "+p.subscription.GetAffSub())
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	q.Add("service_type", "2")
	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.UrlPostbackFsMO+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostbackFsMO + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_FS_MO")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.Msisdn,
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_FS_MO")

	return body, nil
}

/**
 * Delivery Notification (FS)
 */
func (p *Postback) FsDN(status string) ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	if status != "SUCCESS" {
		q.Add("status", "0")
	} else {
		q.Add("status", "1")
	}
	q.Add("statusdesc", strings.ToLower(status))
	q.Add("operator", "321")
	q.Add("msisdn", p.subscription.GetMsisdn())
	q.Add("sdc", "97770")
	q.Add("service", p.subscription.GetCampKeyword()+" "+p.subscription.GetCampSubKeyword())
	q.Add("type", strings.ToLower(p.subscription.GetLatestSubject()))
	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.UrlPostbackFsDN+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostbackFsDN + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_FS_DN")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.Msisdn,
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_FS_DN")

	return body, nil
}

/**
 * Message Originated (PLW)
 */
func (p *Postback) PlwMO() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()
	q := url.Values{}
	q.Add("msisdn", p.subscription.GetMsisdn())

	if p.service.IsMplus() {
		q.Add("id_service", "3131")
		q.Add("operator", "5021")
	}

	q.Add("sms", p.subscription.GetCampKeyword()+" "+p.subscription.GetCampSubKeyword()+" "+p.subscription.GetAffSub())
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	q.Add("service_type", "2")
	q.Add("sdc", "97770")
	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.GetUrlPostbackPlwMO()+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.GetUrlPostbackPlwMO() + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_PLW_MO")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_PLW_MO")

	return body, nil
}

func (p *Postback) PlwMOUnsub() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("msisdn", p.subscription.GetMsisdn())

	if p.service.IsMplus() {
		q.Add("id_service", "3131")
		q.Add("operator", "5021")
	}

	q.Add("sdc", "97770")
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	q.Add("sms", p.subscription.GetLatestKeyword()+" "+p.subscription.GetCampSubKeyword())
	q.Add("service_type", "2")
	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.GetUrlPostbackPlwMO()+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.GetUrlPostbackPlwMO() + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_PLW_MO_UNSUB")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_PLW_MO_UNSUB")

	return body, nil
}

/**
 * Delivery Notification (PLW)
 */
func (p *Postback) PlwDN(status string) ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("msisdn", p.subscription.GetMsisdn())
	if p.service.IsMplus() {
		q.Add("id_service", "3131")
		q.Add("operator", "5021")
	}

	q.Add("trx_id", p.subscription.GetLatestTrxId())
	if status != "SUCCESS" {
		q.Add("status", "0")
	} else {
		q.Add("status", "1")
	}
	q.Add("statusdesc", strings.ToLower(status))
	q.Add("sdc", "97770")
	q.Add("service", p.subscription.GetCampKeyword()+" "+p.subscription.GetCampSubKeyword())
	q.Add("type", strings.ToLower(p.subscription.GetLatestSubject()))
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.GetUrlPostbackPlwDN()+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.GetUrlPostbackPlwDN() + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_PLW_MT")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.Msisdn,
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_PLW_MT")

	return body, nil
}

/**
 * Message Originated (STAR)
 */
func (p *Postback) StarMO() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("cid", p.subscription.GetAffSub())

	req, err := http.NewRequest("GET", p.service.GetUrlPostbackStarMO()+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.UrlPostbackStarMO + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_STAR_MO")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_STAR_MO")

	return body, nil
}

/**
 * Message Originated (MXO)
 */

func (p *Postback) MxoMO() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("msisdn", p.subscription.GetMsisdn())
	// msisdn, id_service, operator, sms, trx_id, service_type, sdc, trx_date

	q.Add("sdc", "97770")
	q.Add("sms", p.subscription.GetCampKeyword()+" "+p.subscription.GetCampSubKeyword()+" "+p.subscription.GetAffSub())
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	q.Add("service_type", "2")
	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.GetUrlPostbackMxoMO()+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.GetUrlPostbackMxoMO() + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_MXO_MO")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_MXO_MO")

	return body, nil
}

func (p *Postback) MxoMOUnsub() ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("msisdn", p.subscription.GetMsisdn())
	q.Add("sdc", "97770")
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	q.Add("sms", p.subscription.GetLatestKeyword()+" "+p.subscription.GetCampSubKeyword())
	q.Add("service_type", "2")

	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.GetUrlPostbackMxoMO()+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  p.subscription.GetMsisdn(),
		"request": p.service.GetUrlPostbackMxoMO() + "?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("POSTBACK_MXO_MO_UNSUB")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_MXO_MO_UNSUB")

	return body, nil
}

/**
 * Delivery Notification (MXO)
 */
func (p *Postback) MxoDN(status string) ([]byte, error) {
	l := p.logger.Init("pb", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	q := url.Values{}
	q.Add("msisdn", p.subscription.GetMsisdn())

	if status != "SUCCESS" {
		q.Add("status", "0")
	} else {
		q.Add("status", "1")
	}

	// msisdn, id_service, operator, trx_id, status, statusdesc, sdc, trx_date
	q.Add("statusdesc", strings.ToLower(status))
	q.Add("sdc", "97770")
	q.Add("service", p.subscription.GetCampKeyword()+" "+p.subscription.GetCampSubKeyword())
	q.Add("type", strings.ToLower(p.subscription.GetLatestSubject()))
	q.Add("trx_id", p.subscription.GetLatestTrxId())
	q.Add("trx_date", time.Now().Format("20060102150405"))

	req, err := http.NewRequest("GET", p.service.GetUrlPostbackMxoDN()+"?"+q.Encode(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{"msisdn": p.subscription.Msisdn, "request": p.service.GetUrlPostbackMxoDN() + "?" + q.Encode(), "trx_id": trxId}).Info("POSTBACK_MXO_DN")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	duration := time.Since(start).Milliseconds()
	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      p.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("POSTBACK_MXO_DN")

	return body, nil
}

func (p *Postback) ManualHit(reqUrl string) ([]byte, error) {
	l := p.logger.Init("pb", true)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	p.logger.Writer(req)
	l.WithFields(logrus.Fields{"request": reqUrl}).Info("POSTBACK_SAM_DN_MANUAL")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	p.logger.Writer(string(body))
	l.WithFields(logrus.Fields{"response": string(body)}).Info("POSTBACK_SAM_DN_MANUAL")

	return body, nil
}
