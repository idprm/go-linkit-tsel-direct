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

	"github.com/idprm/go-linkit-tsel/src/config"
	"github.com/idprm/go-linkit-tsel/src/logger"
	"github.com/sirupsen/logrus"
)

type Arpu struct {
	cfg    *config.Secret
	logger *logger.Logger
}

func NewArpu(
	cfg *config.Secret,
	logger *logger.Logger,
) *Arpu {
	return &Arpu{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *Arpu) UploadCSV(urlTo, fileName string) {
	l := a.logger.Init("csv", true)

	request, err := a.fileUploadRequest(urlTo, "file", fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			l.WithFields(logrus.Fields{"error": err.Error()}).Error("UPLOAD_CSV")
			log.Fatal(err)
		}
		fmt.Println(resp.Header)
		fmt.Println(body)

		defer resp.Body.Close()
		l.WithFields(logrus.Fields{"response": body.String()}).Info("UPLOAD_CSV")
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
	req.SetBasicAuth(a.cfg.Arpu.Username, a.cfg.Arpu.Password)
	req.Header.Add("Bearer", a.cfg.Arpu.Token)

	return req, err
}
