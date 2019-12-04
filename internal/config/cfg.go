package config

import (
	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"sync"
)

const (
	mode        = "O_MODE"
	DevelopMode = "develop"
	TestMode    = "test"
	ReleaseMode = "bjac"
)

var (
	AppName = ""
	ENV     = "develop"
	cfg     *toml.Tree
	_       sync.Once
)

func init() {
	cfgPath := os.Getenv("CONFIG")
	if len(cfgPath) == 0 {
		cfgPath, _ = filepath.Abs(filepath.Dir("../../config/"))
	}

	tfg, err := toml.LoadFile(cfgPath + "/config.toml")
	if err != nil {
		panic(err)
	}

	// app ENV
	switch os.Getenv(mode) {
	case ReleaseMode:
		ENV = ReleaseMode
	case TestMode:
		ENV = TestMode
	default:
		ENV = DevelopMode
	}

	if name := tfg.Get("name"); name != nil {
		AppName = name.(string)
	}

	if tfg.Has(ENV) == false {
		log.Fatal().Caller().Msgf("there is nil %s config", ENV)
	}
	cfg = tfg.Get(ENV).(*toml.Tree)
}

func Ini() *toml.Tree {
	return cfg
}
