package executor

import "time"

type TickerHandlerFunc func(t time.Time)
