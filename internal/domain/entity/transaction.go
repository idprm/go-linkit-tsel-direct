package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	ID             int64     `json:"id,omitempty"`
	TxID           string    `json:"tx_id,omitempty"`
	ServiceID      int       `json:"service_id,omitempty"`
	Service        *Service  `json:",omitempty"`
	Msisdn         string    `json:"msisdn"`
	Channel        string    `json:"channel,omitempty"`
	CampKeyword    string    `json:"camp_keyword,omitempty"`
	CampSubKeyword string    `json:"camp_sub_keyword,omitempty"`
	Adnet          string    `json:"adnet,omitempty"`
	PubID          string    `json:"pub_id,omitempty"`
	AffSub         string    `json:"aff_sub,omitempty"`
	Keyword        string    `json:"keyword,omitempty"`
	PIN            string    `json:"pin,omitempty"`
	Amount         float64   `json:"amount,omitempty"`
	Status         string    `json:"status,omitempty"`
	StatusCode     string    `json:"status_code,omitempty"`
	StatusDetail   string    `json:"status_detail,omitempty"`
	Subject        string    `json:"subject,omitempty"`
	IpAddress      string    `json:"ip_address,omitempty"`
	Payload        string    `json:"payload,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

func (t *Transaction) SetAmount(amount float64) {
	t.Amount = amount
}

func (t *Transaction) SetStatus(status string) {
	t.Status = status
}

func (t *Transaction) SetStatusCode(statusCode string) {
	t.StatusCode = statusCode
}

func (t *Transaction) SetStatusDetail(statusDetail string) {
	t.StatusDetail = statusDetail
}

func (t *Transaction) SetSubject(subject string) {
	t.Subject = subject
}

func (t *Transaction) SetCampKeyword(data string) {
	t.CampKeyword = strings.ToUpper(data)
}

func (t *Transaction) SetCampSubKeyword(data string) {
	t.CampSubKeyword = strings.ToUpper(data)
}

func (t *Transaction) GetAmount() string {
	return strconv.Itoa(int(t.Amount))
}

func (t *Transaction) GetAmountWithSeparator() string {
	return IntComma(int(t.Amount))
}

func IntComma(i int) string {
	if i < 0 {
		return "-" + IntComma(-i)
	}
	if i < 1000 {
		return fmt.Sprintf("%d", i)
	}
	return IntComma(i/1000) + "," + fmt.Sprintf("%03d", i%1000)
}
