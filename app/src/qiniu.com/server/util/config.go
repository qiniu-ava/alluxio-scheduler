package util

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"qiniu.com/server/typo"
)

func GetConfig(p PHASE_TYPE) *typo.Config {
	var confName string
	switch p {
	case PHASE.DEV:
		confName = "config.dev.json"
		break
	case PHASE.PRODUCTION:
		confName = "config.prd.json"
		break
	default:
		confName = ""
		break
	}

	conf := &typo.Config{}
	var loader *confita.Loader
	if confName == "" {
		loader = confita.NewLoader(file.NewBackend("./config/config.default.json"))
	} else {
		loader = confita.NewLoader(file.NewBackend("./config/config.default.json"), file.NewBackend("./config/"+confName))
	}
	loader.Load(context.Background(), conf)

	return conf
}
