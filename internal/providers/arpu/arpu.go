package arpu

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/logger"
	"github.com/idprm/go-linkit-tsel/internal/utils"
	"github.com/sirupsen/logrus"
)

var (
	ARPU_USERNAME string = utils.GetEnv("ARPU_USERNAME")
	ARPU_PASSWORD string = utils.GetEnv("ARPU_PASSWORD")
	ARPU_TOKEN    string = utils.GetEnv("ARPU_TOKEN")
)

type Arpu struct {
	logger *logger.Logger
}

func NewArpu(
	logger *logger.Logger,
) *Arpu {
	return &Arpu{
		logger: logger,
	}
}

func (a *Arpu) UploadCSV(urlTo, fileName string) {
	l := a.logger.Init("csv", true)

	start := time.Now()

	request, err := a.fileUploadRequest(urlTo, "file", fileName)
	if err != nil {
		l.WithFields(logrus.Fields{"error": err.Error()}).Error("UPLOAD_CSV")
		log.Println(err.Error())
	}
	tr := &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       0,
		TLSHandshakeTimeout:   0,
		ResponseHeaderTimeout: 0,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Timeout:   0,
		Transport: tr,
	}
	resp, err := client.Do(request)
	if err != nil {
		l.WithFields(logrus.Fields{"error": err.Error()}).Error("UPLOAD_CSV")
		log.Println(err.Error())
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			l.WithFields(logrus.Fields{"error": err.Error()}).Error("UPLOAD_CSV")
			log.Println(err.Error())
		}
		fmt.Println(resp.Header)
		fmt.Println(body)

		defer resp.Body.Close()
		duration := time.Since(start).Milliseconds()
		l.WithFields(logrus.Fields{
			"duration":    duration,
			"response":    body.String(),
			"status_code": resp.StatusCode,
			"status_text": http.StatusText(resp.StatusCode),
		}).Info("UPLOAD_CSV")
	}
}

func (a *Arpu) fileUploadRequest(uri, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))

	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth(ARPU_USERNAME, ARPU_PASSWORD)
	req.Header.Add("Bearer", ARPU_TOKEN)

	return req, err
}
