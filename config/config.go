package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
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
	MigrationDir string `yaml:"migration-dir"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type JwtConfig struct {
	AccessTokenDuration  string `yaml:"access-token-duration"`
	RefreshTokenDuration string `yaml:"refresh-token-duration"`
	JwtAccessSecret      string `yaml:"jwt-access-secret"`
	JwtRefreshSecret     string `yaml:"jwt-refresh-secret"`
}

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	}
	Redis *RedisConfig
	Db    *DatabaseConfig
	Jwt   *JwtConfig
}

func InitCfg(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
