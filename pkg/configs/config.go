package configs

import (
	"log"
	"math/rand"
	"os"

	"gopkg.in/yaml.v2"
)

var proxyConfig Config
var roundRobin map[string]int

// Config of the reverse proxy
type Config struct {
	Proxy struct {
		Listen   Host      `yaml:"listen"`
		Services []Service `yaml:"services"`
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

// Load parses the config YAML file
func Load(fileName string) (Config, error) {
	configFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&proxyConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer configFile.Close()

	initializeRoundRobin()
	return proxyConfig, nil
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
func ChooseInstance(domain string, instances []Host, isRoundRobin *bool) int {
	// forwards request to an instance in a round-robin strategy if configured
	if *isRoundRobin {
		roundRobin[domain] = (roundRobin[domain] + 1) % len(instances)
		return roundRobin[domain]
	}

	// per default forwards request to an instance randomly
	return rand.Intn(len(instances))
}
