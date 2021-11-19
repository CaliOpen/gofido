package main

import (
	"flag"
	"fmt"
	"github.com/CaliOpen/gofido/config"
	"github.com/CaliOpen/gofido/store"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func readConfig(filename string) (*config.Config, error) {

	conf := &config.Config{}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("Unable to read file ", filename, ": ", filename, err)
		return conf, err
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatal("Unmarshal yaml failed: ", err)
		return &config.Config{}, err
	}
	return conf, nil
}

func main() {
	configfile := flag.String("c", "gofi.yaml", "yaml configuration file")
	flag.Parse()
	config, err := readConfig(*configfile)
	if err != nil {
		fmt.Errorf("Aborting %v", err)
		return
	}

	server := &FidoServer{}
	store := &store.FidoStore{}
	err = store.Initialize(config)
	if err != nil {
		fmt.Errorf("Store initialization fail %v", err)
		return
	}
	err = server.Initialize(config, store)
	if err != nil {
		fmt.Errorf("Server initialization fail %v", err)
		return
	}
	server.Run()
}
