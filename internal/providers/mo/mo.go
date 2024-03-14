package mo

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/idprm/go-linkit-tsel/src/domain/entity"
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
