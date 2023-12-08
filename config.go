package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Conf struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
	Server      string `yaml:"server"`
}

var Config Conf

func init() {
	err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func readConfig() error {
	// open config file
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal("readConfig.open:", err)
	}

	// read config file
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("readConfig.Decode:", err)
	}

	f.Close()
	return nil
}
