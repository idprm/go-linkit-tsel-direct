package entity

import "time"

type Logger struct {
	Type     string    `bson:"type"`
	Msisdn   string    `bson:"msisdn,omitempty"`
	Keyword  string    `bsoh:"keyword"`
	Date     time.Time `bson:"date,omitempty"`
	Request  string    `bson:"request,omitempty"`
	Response string    `bson:"response,omitempty"`
	Message  string    `bson:"message,omitempty"`
}
