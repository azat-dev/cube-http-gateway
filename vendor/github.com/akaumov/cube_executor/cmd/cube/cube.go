package main

import (
	"github.com/akaumov/cube-executor"
	"flag"
	"log"
	"io/ioutil"
	"fmt"
	"encoding/json"
)

func main() {
	var configPath = flag.String("config", "/config.json", "config path")

	flag.Parse()
	log.SetFlags(0)

	configRaw, err := ioutil.ReadFile(*configPath)
	if err != nil {
		fmt.Printf("can't read config file: %v/n", err)
	}

	var config cube_executor.CubeConfig
	err = json.Unmarshal(configRaw, &config)
	if err != nil {
		fmt.Printf("can't parse config file: %v/n", err)
	}

	cube, err := cube_executor.NewCube(config, nil)
	if err != nil {
		log.Fatalf("can't init cube instance: %v/n", err)
	}

	cube.Start()
}
