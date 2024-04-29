package telco

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/idprm/go-linkit-tsel/internal/utils"
	"github.com/idprm/go-linkit-tsel/internal/utils/hash_utils"
	"github.com/idprm/go-linkit-tsel/internal/utils/uuid_utils"
	"github.com/sirupsen/logrus"
)

var (
	TELCO_URL_AUTH string = utils.GetEnv("TELCO_URL_AUTH")
	TELCO_KEY      string = utils.GetEnv("TELCO_KEY")
	TELCO_SECRET   string = utils.GetEnv("TELCO_SECRET")
	TELCO_CPNAME   string = utils.GetEnv("TELCO_CPNAME")
	TELCO_CPID     string = utils.GetEnv("TELCO_CPID")
	TELCO_PWD      string = utils.GetEnv("TELCO_PWD")
	TELCO_SENDER   string = utils.GetEnv("TELCO_SENDER")
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

	req, err := http.NewRequest("GET", t.service.GetUrlTelco()+"/scrt/1/generate.php", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("cp_name", TELCO_CPNAME)
	q.Add("pwd", TELCO_PWD)
	q.Add("programid", t.service.GetProgramId())
	q.Add("sid", t.service.GetSid())

	if t.service.IsGupi() || t.service.IsMplus() {
		q.Add("par", "1")
	}

	req.URL.RawQuery = q.Encode()

	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	strData := TELCO_KEY + TELCO_SECRET + timeStamp

	signature := hash_utils.GetMD5Hash(strData)

	req.Header.Set("api_key", TELCO_KEY)
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
		"request": t.service.GetUrlTelco() + "/scrt/1/generate.php?" + q.Encode(),
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
	l.WithFields(logrus.Fields{"redirect": TELCO_URL_AUTH + "/transaksi/tauthwco?token=" + string(token)}).Info("MT_OPTIN")
	return TELCO_URL_AUTH + "/transaksi/tauthwco?token=" + string(token), string(token), nil
}

func (t *Telco) WebOptInUSSD() (string, error) {
	token, err := t.Token()
	if err != nil {
		return "", err
	}
	return TELCO_URL_AUTH + "/transaksi/konfirmasi/ussd?token=" + string(token), nil
}

func (t *Telco) WebOptInCaptcha() (string, error) {
	token, err := t.Token()
	if err != nil {
		return "", err
	}
	return TELCO_URL_AUTH + "/transaksi/captchawco?token=" + string(token), nil
}

func (t *Telco) SMSbyParam() ([]byte, error) {
	l := t.logger.Init("mt", true)
	//
	start := time.Now()
	trxId := uuid_utils.GenerateTrxId()

	req, err := http.NewRequest(http.MethodGet, t.service.GetUrlTelco()+"/scrt/cp/submitSM.jsp", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("cpid", TELCO_CPID)
	q.Add("sender", TELCO_SENDER)
	q.Add("sms", t.content.GetValue())
	q.Add("pwd", TELCO_PWD)
	q.Add("msisdn", t.subscription.GetMsisdn())
	q.Add("sid", t.service.GetSid())
	q.Add("tid", t.content.GetTid())

	req.URL.RawQuery = q.Encode()

	now := time.Now()
	timeStamp := strconv.Itoa(int(now.Unix()))
	strData := TELCO_KEY + TELCO_SECRET + timeStamp

	signature := hash_utils.GetMD5Hash(strData)

	req.Header.Add("Accept-Charset", "utf-8")
	req.Header.Set("api_key", TELCO_KEY)
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
		"request": t.service.GetUrlTelco() + "/scrt/cp/submitSM.jsp?" + q.Encode(),
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
