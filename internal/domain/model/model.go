package model

type LoggerFormat struct {
	TrxId      string `json:"trx_id"`
	UrlRequest string `json:"url_request"`
	StatusCode string `json:"status_code"`
	StatusText string `json:"status_text"`
	Duration   int    `json:"duration"`
}

type WebResponse struct {
	Error      bool   `json:"error"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IpAddress  string `json:"ip_address,omitempty"`
}

func (m *WebResponse) SetIpAddress(data string) {
	m.IpAddress = data
}
