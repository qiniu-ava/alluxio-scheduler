package util

import "github.com/sirupsen/logrus"

var Logger = logrus.New()

func InitLogger(level logrus.Level) {
	Logger.SetLevel(level)
}
