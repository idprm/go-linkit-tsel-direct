package logger

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	LOG_PATH string = "LOG_PATH"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Init(path string, display bool) *logrus.Logger {
	f, err := os.OpenFile(LOG_PATH+"/"+path+"/"+path+"-"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}

	logger := logrus.New()

	if display {
		logger.SetOutput(io.MultiWriter(os.Stdout, f))
	} else {
		logger.SetOutput(io.MultiWriter(f))
	}

	logger.SetReportCaller(false)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}

func (l *Logger) Writer(data interface{}) {
	//create your file with desired read/write permissions
	f, err := os.OpenFile(LOG_PATH+"/http/log-"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

	log.Println(data)
}
