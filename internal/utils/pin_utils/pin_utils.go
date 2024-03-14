package pin_utils

import (
	"math/rand"
	"strings"
)

func Generate(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func GetLatestMsisdn(msisdn string, limit int) string {
	str := strings.NewReplacer("=", "", "+", "", "/", "")
	message := str.Replace(msisdn)
	return message[len(message)-limit:]
}
