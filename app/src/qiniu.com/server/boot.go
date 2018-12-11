package main

import (
	"time"

	"qiniu.com/server/executor"

	"qiniu.com/server/typo"
	"qiniu.com/server/util"
)

var tickerHandler = TickerHandler{}

func Boot(conf *typo.Config) {
	ticker := time.NewTicker(TickPeriod * time.Second)
	initHandlers(conf)

	go func(th TickerHandler) {
		for t := range ticker.C {
			th.Run(t)
		}
	}(tickerHandler)
	for {
		time.Sleep(10 * time.Second)
	}
}

func initHandlers(conf *typo.Config) {
	tickerHandler.handlers = []executor.TickerHandlerFunc{}
	tickerHandler.RegisterHandler(executor.GetPodHandler(conf))
}

type TickerHandler struct {
	handlers []executor.TickerHandlerFunc
}

func (t *TickerHandler) RegisterHandler(f executor.TickerHandlerFunc) {
	t.handlers = append(t.handlers, f)
	util.Logger.Debugf("running in RegisterHandler, with %d handlers", len(t.handlers))
}

func (t *TickerHandler) Run(tm time.Time) {
	util.Logger.Debugf("running in tickerhandler, with %d handlers", len(t.handlers))
	for index := range t.handlers {
		util.Logger.Debugf("index: %d, handler: %v", index, t.handlers[index])
		t.handlers[index](tm)
	}
}
