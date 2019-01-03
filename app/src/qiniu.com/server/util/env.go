package util

type PHASE_TYPE int

var (
	dev        PHASE_TYPE = 0
	staging    PHASE_TYPE = 1
	alpha      PHASE_TYPE = 2
	beta       PHASE_TYPE = 3
	production PHASE_TYPE = 4
)

type phase struct {
	DEV        PHASE_TYPE
	STAGING    PHASE_TYPE
	ALPHA      PHASE_TYPE
	BETA       PHASE_TYPE
	PRODUCTION PHASE_TYPE
}

var PHASE = phase{
	DEV:        dev,
	STAGING:    staging,
	ALPHA:      alpha,
	BETA:       beta,
	PRODUCTION: production,
}
