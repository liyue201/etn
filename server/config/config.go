package config

import "github.com/jinzhu/configor"

type (
	DbConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		Dbname   string
	}

	Configuration struct {
		Version int
		Port    int
		Static string
		Db      DbConfig
	}
)

var Cfg Configuration

func InitConfig(filePath string) error {
	err := configor.Load(&Cfg, filePath)
	return err
}
