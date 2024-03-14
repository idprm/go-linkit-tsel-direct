package telco

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/idprm/go-linkit-tsel/internal/utils/hash_utils"
	"github.com/idprm/go-linkit-tsel/internal/utils/uuid_utils"
	"github.com/sirupsen/logrus"
)

type Telco struct {
	logger       *logger.Logger
	subscription *entity.Subscription
	service      *entity.Service
	content      *entity.Content
}

func NewTelco(
	logger *logger.Logger,
	subscription *entity.Subscription,
	service *entity.Service,
	content *entity.Content,
) *Telco {
	return &Telco{
		logger:       logger,
		subscription: subscription,
		service:      service,
		content:      content,
	}
}

type ITelco interface {
	Token() ([]byte, error)
	WebOptInOTP() (string, error)
	WebOptInUSSD() (string, error)
	WebOptInCaptcha() (string, error)
	SMSbyParam() ([]byte, error)
}

func (t *Telco) Token() ([]byte, error) {
	l := t.logger.Init("mt", true)

	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	req, err := http.NewRequest("GET", t.cfg.Telco.UrlKey+"/scrt/1/generate.php", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("cp_name", t.cfg.Telco.CpName)
	q.Add("pwd", t.cfg.Telco.Pwd)
	q.Add("programid", t.service.GetProgramId())
	q.Add("sid", t.service.GetSid())

	req.URL.RawQuery = q.Encode()

	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	strData := t.cfg.Telco.Key + t.cfg.Telco.Secret + timeStamp

	signature := hash_utils.GetMD5Hash(strData)

	req.Header.Set("api_key", t.cfg.Telco.Key)
	req.Header.Set("x-signature", signature)

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	t.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"request": t.cfg.Telco.UrlKey + "/scrt/1/generate.php?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("MT_TOKEN")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Milliseconds()
	t.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("MT_TOKEN")

	return body, nil
}

func (t *Telco) WebOptInOTP() (string, string, error) {
	l := t.logger.Init("mt", true)

	token, err := t.Token()
	if err != nil {
		return "", "", err
	}
	l.WithFields(logrus.Fields{"redirect": t.cfg.Telco.UrlAuth + "/transaksi/tauthwco?token=" + string(token)}).Info("MT_OPTIN")
	return t.cfg.Telco.UrlAuth + "/transaksi/tauthwco?token=" + string(token), string(token), nil
}

func (t *Telco) WebOptInUSSD() (string, error) {
	token, err := t.Token()
	if err != nil {
		return "", err
	}
	return t.cfg.Telco.UrlAuth + "/transaksi/konfirmasi/ussd?token=" + string(token), nil
}

func (t *Telco) WebOptInCaptcha() (string, error) {
	token, err := t.Token()
	if err != nil {
		return "", err
	}
	return t.cfg.Telco.UrlAuth + "/transaksi/captchawco?token=" + string(token), nil
}

func (t *Telco) SMSbyParam() ([]byte, error) {
	l := t.logger.Init("mt", true)
	//
	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	req, err := http.NewRequest(http.MethodGet, t.cfg.Telco.UrlKey+"/scrt/cp/submitSM.jsp", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("cpid", t.cfg.Telco.CpId)
	q.Add("sender", t.cfg.Telco.Sender)
	q.Add("sms", t.content.GetValue())
	q.Add("pwd", t.cfg.Telco.Pwd)
	q.Add("msisdn", t.subscription.GetMsisdn())
	q.Add("sid", t.service.GetSid())
	q.Add("tid", t.content.GetTid())

	req.URL.RawQuery = q.Encode()

	now := time.Now()
	timeStamp := strconv.Itoa(int(now.Unix()))
	strData := t.cfg.Telco.Key + t.cfg.Telco.Secret + timeStamp

	signature := hash_utils.GetMD5Hash(strData)

	req.Header.Add("Accept-Charset", "utf-8")
	req.Header.Set("api_key", t.cfg.Telco.Key)
	req.Header.Set("x-signature", signature)

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	t.logger.Writer(req)
	l.WithFields(logrus.Fields{
		"msisdn":  t.subscription.GetMsisdn(),
		"request": t.cfg.Telco.UrlKey + "/scrt/cp/submitSM.jsp?" + q.Encode(),
		"trx_id":  trxId,
	}).Info("MT_SMS")

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
	t.logger.Writer(string(body))
	l.WithFields(logrus.Fields{
		"msisdn":      t.subscription.GetMsisdn(),
		"response":    string(body),
		"trx_id":      trxId,
		"duration":    duration,
		"status_code": resp.StatusCode,
		"status_text": http.StatusText(resp.StatusCode),
	}).Info("MT_SMS")

	return body, nil
}
