// +build debug

package main

import (
	"github.com/sirupsen/logrus"
	"qiniu.com/server/util"
)

const TickPeriod = 5

func main() {
	util.InitLogger(logrus.DebugLevel)
	conf := util.GetConfig(util.PHASE.DEV)
	util.Logger.Infof("config: %t", conf)
	Boot(conf)
}
