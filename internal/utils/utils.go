package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

func GenerateTrxId() string {
	id := uuid.New()
	return id.String()
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Panicf("Error %v", key)
	}
	return value
}

func GetLatestMsisdn(msisdn string, limit int) string {
	str := strings.NewReplacer("=", "", "+", "", "/", "")
	message := str.Replace(msisdn)
	return message[len(message)-limit:]
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
