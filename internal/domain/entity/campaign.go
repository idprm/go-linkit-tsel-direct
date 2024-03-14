package entity

type Campaign struct {
	ID        int `json:"id"`
	ServiceID int `json:"service_id"`
	Service   *Service
	Adnet     string `json:"adnet"`
	LimitMo   int    `json:"limit_mo"`
	TotalMo   int    `json:"total_mo"`
}
