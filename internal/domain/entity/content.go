package entity

import (
	"strings"
)

type Content struct {
	ID        int `json:"id"`
	ServiceID int `json:"service_id"`
	Service   *Service
	Name      string `json:"name"`
	Value     string `json:"value"`
	Tid       string `json:"tid"`
}

func (c *Content) GetName() string {
	return c.Name
}

func (c *Content) GetValue() string {
	return c.Value
}

func (c *Content) GetTid() string {
	return c.Tid
}

func (c *Content) SetPIN(pin string) {
	replacer := strings.NewReplacer("@pin", pin)
	c.Value = replacer.Replace(c.Value)
}

func (c *Content) SetLinkPortalMainPlus(val string) {
	replacer := strings.NewReplacer("https://tsel.mainplus.mobi/", val)
	c.Value = replacer.Replace(c.Value)
}
