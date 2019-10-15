package cmd

import (
	"io"
	"log"
	"net/http"

	"../pkg/configs"
)

// handler sends requests to a service instance and forwards response back to origin
func handler(w http.ResponseWriter, r *http.Request) {
	// get the map of available downstream services
	services, err := configs.GetServices()
	if err != nil {
		log.Fatal(err)
	}

	// choose a instance to forward requests based on load balancing strategy
	instance := configs.ChooseInstance(r.Host)

	// make request to chosen instance
	resp, err := http.Get(services[r.Host][instance].Address + ":" + services[r.Host][instance].Port)
	if err != nil {
		log.Fatal(err)
	}

	// receive and copy response body to forward back to origin
	io.Copy(w, resp.Body)
	resp.Body.Close()
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
	// Read config YAML file
	proxyConfig, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Start proxy
	if err := startProxy(&proxyConfig); err != nil {
		log.Fatal(err)
	}
}
