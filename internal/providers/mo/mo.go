package mo

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/idprm/go-linkit-tsel/internal/domain/model"
)

func HitMO(r entity.ReqSub) ([]byte, error) {

	req, err := http.NewRequest("GET", "https://linkit.exmp.fun/mo", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("sms", r.Sms)
	q.Add("adn", r.Adn)
	q.Add("msisdn", r.Msisdn)
	q.Add("ip_address", r.IpAddress)

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		MaxIdleConns:       20,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return []byte(body), nil
}

func HitPostback(r model.RePostBackResponse) ([]byte, error) {

	req, err := http.NewRequest("GET", "http://kbtools.net/id-linkittisel.php", nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("partner", "linkittisel")
	q.Add("px", r.Verify.AffSub)
	q.Add("serv_id", r.Verify.Service+" "+r.Verify.CampSubKeyword)
	q.Add("token", r.Verify.Token)
	q.Add("msisdn", r.Request.Msisdn)
	q.Add("trxid", r.Verify.TrxId)
	q.Add("time", time.Now().String())

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")

	log.Println(req)

	tr := &http.Transport{
		MaxIdleConns:       20,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	log.Println(string(body))

	return []byte(body), nil
}
