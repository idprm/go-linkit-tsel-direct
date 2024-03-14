package model

type LoggerFormat struct {
	TrxId      string `json:"trx_id"`
	UrlRequest string `json:"url_request"`
	StatusCode string `json:"status_code"`
	StatusText string `json:"status_text"`
	Duration   int    `json:"duration"`
}
