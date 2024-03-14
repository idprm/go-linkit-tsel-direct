package logger

import (
	"io"
	"os"
	"time"

	"github.com/idprm/go-linkit-tsel/src/config"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	cfg *config.Secret
}

func NewLogger(cfg *config.Secret) *Logger {
	return &Logger{
		cfg: cfg,
	}
}

func (l *Logger) Init(path string, display bool) *logrus.Logger {
	f, err := os.OpenFile(l.cfg.Log.Path+"/"+path+"/"+path+"-"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
