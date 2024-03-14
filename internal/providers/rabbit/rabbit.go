package rabbit

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/utils/auth_utils"
)

type RabbitMQ struct {
}

func NewRabbitMQ() *RabbitMQ {
	return &RabbitMQ{}
}

func (p *RabbitMQ) Queue(name string) ([]byte, error) {
	req, err := http.NewRequest("GET", p.cfg.GetUrlRabbitMq()+name, nil)
	req.Header.Add("Authorization", "Basic "+auth_utils.BasicAuth(p.cfg.Rmq.User, p.cfg.Rmq.Pass))

	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
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

	return body, nil
}

func (p *RabbitMQ) Purge(name string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", p.cfg.GetUrlRabbitMq()+name+"/contents", nil)
	req.Header.Add("Authorization", "Basic "+auth_utils.BasicAuth(p.cfg.Rmq.User, p.cfg.Rmq.Pass))

	if err != nil {
		return nil, errors.New(err.Error())
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
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

	return body, nil
}
