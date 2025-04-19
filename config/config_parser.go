package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Cfg struct {
	Env          string      `yaml:"env"`
	Database_url string      `yaml:"database_url"`
	Server       server_data `yaml:"server"`
}
type server_data struct {
	Timeout       time.Duration `yaml:"timeout"`
	TokenTTL      time.Duration `yaml:"token_ttl"`
	Url_accepter  string        `yaml:"url_accepter"`
	Mail_sender   string        `yaml:"mail_sender"`
	Mail_password string        `yaml:"mail_password"`
	Idletimout    time.Duration `yaml:"idletimeout"`
	Address       string        `yaml:"address"`
	Port          int           `yaml:"port"`
}

func Cfg_parser() Cfg {
	var config_path string
	flag.StringVar(&config_path, "config", "", "config file path")
	flag.Parse()
	if config_path == "" {
		config_path = os.Getenv("CONFIG_PATH")
		_, err := os.Stat(config_path)
		if err != nil {
			log.Fatal("config file path err:", err)
		}
		if os.IsNotExist(err) {
			log.Fatal("config file path err:", err)
		}
	}
	var cfg Cfg
	err := cleanenv.ReadConfig(config_path, &cfg)
	if err != nil {
		log.Fatal("read config err:", err)
	}
	return cfg
}
