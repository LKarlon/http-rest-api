package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/LKarlon/http-rest-api.git/internal/app/apiserver"
	"log"
)

var(
	configPath string
)

func init(){
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()
	config :=apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	s := apiserver.New(config)
	go s.Worker()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
