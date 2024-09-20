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

type RePostBackResponse struct {
	Request struct {
		Msisdn string `json:"msisdn"`
	} `json:"request"`
	Verify struct {
		TrxId          string `json:"tx_id"`
		Token          string `json:"token"`
		Service        string `json:"service"`
		Adnet          string `json:"adnet"`
		PubId          string `json:"pub_id"`
		AffSub         string `json:"aff_sub"`
		CampSubKeyword string `json:"camp_sub_keyword"`
	} `json:"verify"`
}
