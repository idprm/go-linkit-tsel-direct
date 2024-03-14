package entity

import "time"

type History struct {
	ID        int64 `json:"id"`
	ServiceID int   `json:"service_id"`
	Service   *Service
	Msisdn    string    `json:"msisdn"`
	Channel   string    `json:"channel"`
	Adnet     string    `json:"adnet"`
	Keyword   string    `json:"keyword"`
	Subject   string    `json:"subject"`
	Status    string    `json:"status"`
	IpAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
