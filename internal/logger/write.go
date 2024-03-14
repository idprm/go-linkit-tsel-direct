package logger

import (
	"log"
	"os"
	"time"
)

func (l *Logger) Writer(data interface{}) {
	//create your file with desired read/write permissions
	f, err := os.OpenFile(l.cfg.Log.Path+"/http/log-"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

	log.Println(data)
}
