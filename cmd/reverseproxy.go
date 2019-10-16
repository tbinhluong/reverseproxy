package cmd

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/tbinhluong/reverseproxy/pkg/configs"
)

var (
	configFile = kingpin.Flag("config.file", "Path of configuration YAML file.").Default("config.yaml").String()
	roundRobin = kingpin.Flag("roundrobin", "Enable round-robin as load balancing strategy, otherwise randomly").Default("false").Bool()
)

// handler sends requests to a service instance and forwards response back to origin
func handler(res http.ResponseWriter, req *http.Request) {
	// get the map of available downstream services
	services, err := configs.GetServices()
	if err != nil {
		log.Fatal(err)
	}

	// choose a instance to forward requests based on load balancing strategy
	instanceID := configs.ChooseInstance(req.Host, services[req.Host], roundRobin)
	instanceURL := services[req.Host][instanceID].Address + ":" + services[req.Host][instanceID].Port
	instance, _ := url.Parse(instanceURL)

	// make request to chosen instance
	proxy := httputil.NewSingleHostReverseProxy(instance)
	proxy.ServeHTTP(res, req)

	/*
		resp, err := http.Get()
		if err != nil {
			log.Fatal(err)
		}

		// receive and copy response body to forward back to origin
		io.Copy(w, resp.Body)
		resp.Body.Close()
	*/
}

// startProxy starts the proxy server on configured host and port
func startProxy(proxyConfig *configs.Config) error {
	// Start the proxy and listen on configured address & port
	server := (*proxyConfig).Proxy.Listen.Address + ":" + (*proxyConfig).Proxy.Listen.Port
	http.HandleFunc("/", handler)

	return http.ListenAndServe(server, nil)
}

// main function
func Execute() {
	// Parse command line args to get path of config and load belancing strategy
	kingpin.New(filepath.Base(os.Args[0]), "A simple reverse proxy")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Read config YAML file
	proxyConfig, err := configs.Load(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Start proxy
	if err := startProxy(&proxyConfig); err != nil {
		log.Fatal(err)
	}
}
