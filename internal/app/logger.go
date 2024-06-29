package app

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func SetLogger(level string) {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrusLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	})

	logrus.SetOutput(os.Stdout)
}
