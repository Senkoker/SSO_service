package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type ClientParser struct {
	Addres       string        `yaml:"addres"`
	Port         int           `yaml:"port"`
	Database_url string        `yaml:"database_url"`
	TokenTtl     time.Duration `yaml:"token_ttl"`
	Idletimeout  time.Duration `yaml:"idletimeout"`
}

func ClientConfigparser() ClientParser {
	var configpath string
	flag.StringVar(&configpath, "config", "", "config path")
	flag.Parse()
	configpath = "C:\\Golang_social_project\\GRPC_Service_sso\\config\\test_config.yaml"
	if configpath == "" {
		configpath = os.Getenv("CLIENT_CONFIG_PATH")
		_, err := os.Stat(configpath)
		if os.IsNotExist(err) {
			log.Fatalln("config file not exist")
		}
	}
	var config ClientParser
	err := cleanenv.ReadConfig(configpath, &config)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}
