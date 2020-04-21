package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Mongo struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database""`
		Url      string `yaml:"url"`
		Port     string `yaml:"port""`
	} `yaml:"mongo"`
	Redis struct {
		Url      string `yaml:"url"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redis"`
}

var Configer *Config

func Load() {
	Configer = new(Config)
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, Configer)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
