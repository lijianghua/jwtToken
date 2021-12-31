package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type DatabaseConfig struct {
	Type         string `yaml:"type"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"databasename"`
	MaxIdleConns int    `yaml:"max-idle-conns"`
	MaxOpenConns int    `yaml:"max-open-conns"`
}

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	}
	Redis struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	Db  DatabaseConfig
	Jwt struct {
		AccessTokenDuration  string `yaml:"access-token-duration"`
		RefreshTokenDuration string `yaml:"refresh-token-duration"`
		JwtAccessSecret      string `yaml:"jwt-access-secret"`
		JwtRefreshSecret     string `yaml:"jwt-refresh-secret"`
	}
}

var Cfg Config

func InitCfg(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln("read config file error: " + err.Error())
	}
	err = yaml.Unmarshal(data, &Cfg)
	if err != nil {
		log.Fatalln("config file unmarshal error: " + err.Error())
	}
}
