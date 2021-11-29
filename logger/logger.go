package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var once sync.Once

type AppLogger struct {
	Logger *logrus.Logger
}

var instance AppLogger

func GetLogger() *logrus.Logger {
	once.Do(func() {
		log := logrus.New()
		log.SetReportCaller(true)

		var filename = "logs/app.log"
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		} else {
			log.SetOutput(f)
		}

		log.SetLevel(logrus.DebugLevel)
		formatter := logrus.JSONFormatter{TimestampFormat: "02-01-2006 15:04:05", PrettyPrint: true}
		log.SetFormatter(&formatter)

		instance = AppLogger{Logger: log}
		log.Debug("GetLogger() is called ONCE!") // this is to check singleton functionality ;)
	})

	return instance.Logger
}
