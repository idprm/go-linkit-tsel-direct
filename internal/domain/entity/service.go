package entity

import "strconv"

type Service struct {
	ID                  int     `json:"id"`
	Category            string  `json:"category"`
	Code                string  `json:"code"`
	Name                string  `json:"name"`
	Package             string  `json:"package"`
	Price               float64 `json:"price"`
	ProgramId           string  `json:"program_id"`
	Sid                 string  `json:"sid"`
	RenewalDay          int     `json:"renewal_day"`
	TrialDay            int     `json:"trial_day"`
	UrlTelco            string  `json:"url_telco"`
	UrlPortal           string  `json:"url_portal"`
	UrlCallback         string  `json:"url_callback"`
	UrlNotifSub         string  `json:"url_notif_sub"`
	UrlNotifUnsub       string  `json:"url_notif_unsub"`
	UrlNotifRenewal     string  `json:"url_notif_renewal"`
	UrlPostback         string  `json:"url_postback"`
	UrlPostbackBillable string  `json:"url_postback_billable"`
	UrlPostbackSamMO    string  `json:"url_postback_sam_mo"`
	UrlPostbackSamDN    string  `json:"url_postback_sam_dn"`
	UrlPostbackYlcMO    string  `json:"url_postback_ylc_mo"`
	UrlPostbackYlcMT    string  `json:"url_postback_ylc_mt"`
	UrlPostbackFsMO     string  `json:"url_postback_fs_mo"`
	UrlPostbackFsDN     string  `json:"url_postback_fs_dn"`
	UrlPostbackPlwMO    string  `json:"url_postback_plw_mo"`
	UrlPostbackPlwDN    string  `json:"url_postback_plw_dn"`
	UrlPostbackStarMO   string  `json:"url_postback_star_mo"`
	UrlPostbackStarDN   string  `json:"url_postback_star_dn"`
	UrlPostbackMxoMO    string  `json:"url_postback_mxo_mo"`
	UrlPostbackMxoDN    string  `json:"url_postback_mxo_dn"`
	UrlPostbackStarsMO  string  `json:"url_postback_stars_mo"`
	UrlPostbackUntMO    string  `json:"url_postback_unt_mo"`
	UrlPostbackUntDN    string  `json:"url_postback_unt_dn"`
}

func (s *Service) GetID() int {
	return s.ID
}

func (s *Service) GetCategory() string {
	return s.Category
}

func (s *Service) GetCode() string {
	return s.Code
}

func (s *Service) GetName() string {
	return s.Name
}

func (s *Service) GetPackage() string {
	return strconv.Itoa(s.RenewalDay)
}

func (s *Service) GetPrice() float64 {
	return s.Price
}

func (s *Service) GetProgramId() string {
	return s.ProgramId
}

func (s *Service) GetSid() string {
	return s.Sid
}

func (s *Service) GetRenewalDay() int {
	return s.RenewalDay
}

func (s *Service) GetTrialDay() int {
	return s.TrialDay
}

func (s *Service) GetUrlTelco() string {
	return s.UrlTelco
}

func (s *Service) GetUrlPortal() string {
	return s.UrlPortal
}

func (s *Service) GetUrlCallback() string {
	return s.UrlCallback
}

func (s *Service) GetUrlNotifSub() string {
	return s.UrlNotifSub
}

func (s *Service) GetUrlNotifUnsub() string {
	return s.UrlNotifUnsub
}

func (s *Service) GetUrlNotifRenewal() string {
	return s.UrlNotifRenewal
}

func (s *Service) GetUrlPostback() string {
	return s.UrlPostback
}

func (s *Service) GetUrlPostbackBillable() string {
	return s.UrlPostbackBillable
}

func (s *Service) GetUrlPostbackSamMO() string {
	return s.UrlPostbackSamMO
}

func (s *Service) GetUrlPostbackSamDN() string {
	return s.UrlPostbackSamDN
}

func (s *Service) GetUrlPostbackYlcMO() string {
	return s.UrlPostbackYlcMO
}

func (s *Service) GetUrlPostbackYlcMT() string {
	return s.UrlPostbackYlcMT
}

func (s *Service) GetUrlPostbackFsMO() string {
	return s.UrlPostbackFsMO
}

func (s *Service) GetUrlPostbackFsDN() string {
	return s.UrlPostbackFsDN
}

func (s *Service) GetUrlPostbackPlwMO() string {
	return s.UrlPostbackPlwMO
}

func (s *Service) GetUrlPostbackPlwDN() string {
	return s.UrlPostbackPlwDN
}

func (s *Service) GetUrlPostbackStarMO() string {
	return s.UrlPostbackStarMO
}

func (s *Service) GetUrlPostbackStarDN() string {
	return s.UrlPostbackStarDN
}

func (s *Service) GetUrlPostbackMxoMO() string {
	return s.UrlPostbackMxoMO
}

func (s *Service) GetUrlPostbackMxoDN() string {
	return s.UrlPostbackMxoDN
}

func (s *Service) GetUrlPostbackStarsMO() string {
	return s.UrlPostbackStarsMO
}

func (s *Service) GetUrlPostbackUntMO() string {
	return s.UrlPostbackUntMO
}

func (s *Service) GetUrlPostbackUntDN() string {
	return s.UrlPostbackUntDN
}

func (s *Service) IsCloudplay() bool {
	return s.GetCategory() == "CLOUDPLAY"
}

func (s *Service) IsGalays() bool {
	return s.GetCategory() == "GALAYS"
}

func (s *Service) IsGupi() bool {
	return s.GetCategory() == "GUPI"
}

func (s *Service) IsMplus() bool {
	return s.GetCategory() == "MPLUS"
}
