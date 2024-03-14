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
	req.SetBasicAuth(ARPU_USERNAME, ARPU_PASSWORD)
	req.Header.Add("Bearer", ARPU_TOKEN)

	return req, err
}
