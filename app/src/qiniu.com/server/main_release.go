// +build !debug

package main

import (
	"github.com/sirupsen/logrus"
	"qiniu.com/server/util"
)

const TickPeriod = 300

func main() {
	util.InitLogger(logrus.InfoLevel)
	conf := util.GetConfig(util.PHASE.PRODUCTION)
	util.Logger.Infof("config: %t", conf)
	Boot(conf)
}
