package config

import (
	"Graduation-Project/log"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	MysqlDsn string // MySQL数据源名称，用于数据库连接
}

var Conf Config

func (conf *Config) InitConfig() {
	yamlFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		log.Logger.Errorf("yamlFile.Get err%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Logger.Errorf("Unmarshal: %v", err)
	}
	Conf = *conf
	return
}
