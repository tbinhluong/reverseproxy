package configs

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var proxyConfig Config
var roundRobin map[string]int

// Config of the reverse proxy
type Config struct {
	Proxy struct {
		Listen   Host      `yaml:"listen"`
		Services []Service `yaml:"services"`
		/*
			Services []struct {
				Name   string    `yaml:"name"`
				Domain string    `yaml:"domain"`
				Hosts  []Address `yaml:"hosts"`
			} `yaml:"services"`
		*/
	} `yaml:"proxy"`
}

// Service config in yaml
type Service struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	Hosts  []Host `yaml:"hosts"`
}

// Host config in yaml
type Host struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

// Load parses the config YAML file
func Load() (Config, error) {
	fileName := configFileName()

	configFile, err := os.Open(fileName)
	if err != nil {
		processError(err)
	}

	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&proxyConfig)
	if err != nil {
		processError(err)
	}
	defer configFile.Close()

	initializeRoundRobin()
	return proxyConfig, nil
}

// return the config YAML file
func configFileName() string {
	return os.Args[1]
}

func initializeRoundRobin() {
	roundRobin = make(map[string]int)
	for _, service := range proxyConfig.Proxy.Services {
		roundRobin[service.Domain] = 0
	}
}

// GetServices converts the array of server to a map of domain --> list of instances which serve requests of that domain
func GetServices() (map[string][]Host, error) {
	var services = make(map[string][]Host)

	for _, service := range proxyConfig.Proxy.Services {
		services[service.Domain] = service.Hosts
	}

	return services, nil
}

// ChooseInstance chooses an instance to forward requests based on load balancing strategy
func ChooseInstance(domain string, instances []Host) int {
	if len(os.Args) == 2 && strings.ToUpper(os.Args[2]) == "ROUND-ROBIN" {
		roundRobin[domain] = (roundRobin[domain] + 1) % len(instances)
		return roundRobin[domain]
	}

	return rand.Intn(len(instances)) - 1
}
