package main

import (
	"os"
	"path/filepath"

	"github.com/xrfang/logging/v2"
	"gopkg.in/yaml.v3"
)

type (
	config struct {
		Logging struct {
			Level logging.LogLevel `yaml:"level"`
			Path  string           `yaml:"path"`
			Split int              `yaml:"split"`
			Keep  int              `yaml:"keep"`
		} `yaml:"logging"`
		Redirects []RedirectItem `yaml:"redirects"`
		confDir   string
	}
)

func (c config) absPath(p string) (ap string) {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Clean(filepath.Join(c.confDir, p))
}

var (
	L  logging.Logger
	cf config
)

func loadConfig(fn string) {
	if fn == "" {
		panic("missing config file")
	}
	cf.Logging.Path = "logs"
	cf.Logging.Split = 10 * 1024 * 1024 //每个LOG文件10兆字节
	cf.Logging.Keep = 10
	fp, err := filepath.Abs(fn)
	assert(err)
	cf.confDir = filepath.Dir(fp)
	f, err := os.Open(fp)
	assert(err)
	defer f.Close()
	assert(yaml.NewDecoder(f).Decode(&cf))
	cf.Logging.Path = cf.absPath(cf.Logging.Path)
}
