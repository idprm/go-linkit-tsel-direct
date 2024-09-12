package entity

import "strings"

type Verify struct {
	TxId           string `json:"tx_id,omitempty"`
	Token          string `json:"token,omitempty"`
	Service        string `json:"service,omitempty"`
	Adnet          string `json:"adnet,omitempty"`
	PubID          string `json:"pub_id,omitempty"`
	AffSub         string `json:"aff_sub,omitempty"`
	CampKeyword    string `json:"camp_keyword,omitempty"`
	CampSubKeyword string `json:"camp_sub_keyword,omitempty"`
	Browser        string `json:"browser,omitempty"`
	OS             string `json:"os,omitempty"`
	Device         string `json:"device,omitempty"`
	IpAddress      string `json:"ip_address,omitempty"`
	IsBillable     bool   `json:"is_billable,omitempty"`
	IsCampTool     bool   `json:"is_camptool,omitempty"`
}

func (v *Verify) GetTxId() string {
	return v.TxId
}

func (v *Verify) GetToken() string {
	return v.Token
}

func (v *Verify) GetService() string {
	return v.Service
}

func (v *Verify) GetAdnet() string {
	return v.Adnet
}

func (v *Verify) GetPubId() string {
	return v.PubID
}

func (v *Verify) GetAffSub() string {
	return v.AffSub
}

func (v *Verify) GetCampKeyword() string {
	return strings.ToUpper(v.CampKeyword)
}

func (v *Verify) GetCampSubKeyword() string {
	return strings.ToUpper(v.CampSubKeyword)
}

func (v *Verify) GetBrowser() string {
	return v.Browser
}

func (v *Verify) GetOS() string {
	return v.OS
}

func (v *Verify) GetDevice() string {
	return v.Device
}

func (v *Verify) GetIpAddress() string {
	return v.IpAddress
}

func (v *Verify) GetIsBillable() bool {
	return v.IsBillable
}

func (v *Verify) GetIsCampTool() bool {
	return v.IsCampTool
}

func (v *Verify) SetCampKeyword(keyword string) {
	v.CampKeyword = strings.ToUpper(keyword)
}

func (v *Verify) SetCampSubKeyword(subkey string) {
	v.CampSubKeyword = strings.ToUpper(subkey)
}

func (v *Verify) IsCampKeyword() bool {
	return v.CampKeyword != ""
}

func (v *Verify) IsSam() bool {
	return strings.ToUpper(v.CampSubKeyword) == "SAM"
}

func (v *Verify) IsYlc() bool {
	return strings.ToUpper(v.CampSubKeyword) == "YLC" || strings.ToUpper(v.CampSubKeyword) == "YL2"
}

func (v *Verify) IsBng() bool {
	return strings.ToUpper(v.CampSubKeyword) == "BNG"
}

func (v *Verify) IsFs() bool {
	return strings.ToUpper(v.CampSubKeyword) == "FS"
}

func (v *Verify) IsRdr() bool {
	return strings.ToUpper(v.CampSubKeyword) == "RDR"
}

func (v *Verify) IsV2Test() bool {
	return strings.ToUpper(v.CampSubKeyword) == "V2TEST"
}

func (v *Verify) IsPlw() bool {
	return strings.ToUpper(v.CampSubKeyword) == "PLW"
}

func (v *Verify) IsStar() bool {
	return strings.ToUpper(v.CampSubKeyword) == "STAR"
}

func (v *Verify) IsMxo() bool {
	return strings.ToUpper(v.CampSubKeyword) == "MXO"
}

func (v *Verify) IsStars() bool {
	return strings.ToUpper(v.CampSubKeyword) == "STARS"
}

func (v *Verify) IsUnt() bool {
	return strings.ToUpper(v.CampSubKeyword) == "UNT"
}
