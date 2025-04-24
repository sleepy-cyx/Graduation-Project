package config

import (
	"Graduation-Project/log"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	MysqlDsn       string // MySQL数据源名称，用于数据库连接
	DeepSeekAPIKey string
}

var Conf Config

func InitConfig() {
	yamlFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		log.Logger.Errorf("yamlFile.Get err%v ", err)
	}
	conf := Config{}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Logger.Errorf("Unmarshal: %v", err)
	}
	Conf = conf
	return
}
func GetConfig() Config {
	InitConfig()
	return Conf
}
