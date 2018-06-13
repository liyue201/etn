package config

import (
	"github.com/jinzhu/configor"
)

type (
	Configuration struct {
		Version int
	}
)

var Cfg Configuration

func InitConfig(filePath string) error {
	err := configor.Load(&Cfg, filePath)
	return err
}
