# Copyright 2019 Binh Luong
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

package configs

import (
	"errors"
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

	if len(proxyConfig.Proxy.Services) == 0 {
		return nil, errors.New("No downstream services specified")
	}

	for _, service := range proxyConfig.Proxy.Services {
		services[service.Domain] = service.Hosts
	}

	return services, nil
}

// ChooseInstance chooses an instance to forward requests based on load balancing strategy
func ChooseInstance(domain string, instances []Host, isRoundRobin *bool) (int, error) {
	if len(instances) == 0 {
		return 0, errors.New("Downstream service not found")
	}

	// forwards request to an instance in a round-robin strategy if configured
	if *isRoundRobin {
		roundRobin[domain] = (roundRobin[domain] + 1) % len(instances)
		return roundRobin[domain], nil
	}

	// per default forwards request to an instance randomly
	return rand.Intn(len(instances)), nil
}
