package response_utils

import "strings"

func ParseStatusCode(code string) string {
	switch code {
	case "0:1":
		return "Default error code"
	case "0:2":
		return "MT rejected due to storage partition is full"
	case "1":
		return "Success"
	case "2":
		return "Authentication failed (binding failed)"
	case "3:":
		return "Charging failed"
	case "3:101":
		return "Charging timeout"
	case "3:105":
		return "Invalid MSISDN (recipient)"
	case "3:3:105":
		return "MT is rejected due to MSISDN is blacklist"
	case "3:3:21":
		return "Not enough credit"
	case "3:3:27":
		return "Number is out of Active period (Grace period)"
	case "4:1":
		return "Invalid shortcode (sender)"
	case "4:2:":
		return "Mandatory parameter is missing"
	case "4:3":
		return "MT rejected due to long message restriction"
	case "4:4:1":
		return "Multiple tariff is not allowed, but “tid” parameter is provided by CP"
	case "4:4:2":
		return "The provided “tid” by CP is not allowed"
	case "5:1":
		return "MT rejected due to subscription quota is finished"
	case "5:2":
		return "MT rejected due to subscriber doesn't have this subscription"
	case "5:3":
		return "MT rejected due to subscription is disabled"
	case "5:4":
		return "Throttling error"
	case "5:997":
		return "Invalid trx_id"
	case "6":
		return "MT rejected due to quarantine"
	case "7":
		return "Error XML"
	default:
		return "No description"
	}
}

func ParseStatus(code string) bool {
	return strings.HasPrefix(code, "1")
}

func ParseChannel(sms string) string {
	if strings.Contains(strings.ToUpper(sms), "TOKEN=") {
		return "WAP"
	}
	if strings.Contains(strings.ToUpper(sms), "TOKEN%3D") {
		return "WAP"
	}
	return "SMS"
}

func ParseToken(sms string) string {
	message := strings.ToUpper(sms)
	i := strings.LastIndex(message, " ")
	if i > -1 {
		keyword := message[i+1:]
		if strings.Contains(message, "TOKEN=") {
			token := len("TOKEN=")
			return keyword[token:]
		}
		if strings.Contains(message, "TOKEN%3D") {
			token := len("TOKEN%3D")
			return keyword[token:]
		}
	}
	return ""
}

func IsPurge(code string) bool {
	// 3:105 = Invalid MSISDN (recipient)
	// 3:3:105 = MT is rejected due to MSISDN is blacklist
	// 5:2 = MT rejected due to subscriber doesn't have this subscription
	return code == "3:105" || code == "3:3:105" || code == "5:2"
}

func IsInsuff(code string) bool {
	// 3:3:21 = Not enough credit
	return code == "3:3:21"
}
