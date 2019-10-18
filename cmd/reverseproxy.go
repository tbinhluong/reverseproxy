package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/tbinhluong/reverseproxy/pkg/configs"
)

const scheme = "http"

var (
	configFile = kingpin.Flag("config.file", "Path of configuration YAML file.").Default("config.yml").String()
	roundRobin = kingpin.Flag("roundrobin", "Enable round-robin as load balancing strategy, otherwise randomly").Default("false").Bool()
)

// handler sends requests to a service instance and forwards response back to origin
func handler(res http.ResponseWriter, req *http.Request) {

	// get the map of available downstream services
	services, err := configs.GetServices()
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// choose a instance to forward requests based on load balancing strategy
	instanceID, err := configs.ChooseInstance(req.Host, services[req.Host], roundRobin)
	if err != nil {
		// check if it's a health check request
		if isHealthCheck(req) {
			log.Println("Health check")
			return
		}
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	instanceURL := scheme + "://" + services[req.Host][instanceID].Address + ":" + services[req.Host][instanceID].Port
	instance, _ := url.Parse(instanceURL)

	// make request to chosen instance
	log.Println("Forwarding requests", req.RequestURI, "to", req.Host, "on", instanceURL)
	proxy := httputil.NewSingleHostReverseProxy(instance)
	proxy.ServeHTTP(res, req)
}

func isHealthCheck(req *http.Request) bool {
	if req.RequestURI == "/healthz" && req.Method == "GET" {
		return true
	}
	return false
}

// startProxy starts the proxy server on configured host and port
func startProxy(proxyConfig *configs.Config) error {
	// Start the proxy and listen on configured address & port
	server, err := getProxyServer(proxyConfig)
	if err != nil {
		return err
	}
	log.Println("Start listening to HTTP requests on", server)

	http.HandleFunc("/", handler)

	// in order to run as pod in k8s cluster, it should listen on any interface
	listenAllInterfaces := strings.Replace(server, (*proxyConfig).Proxy.Listen.Address, "", -1)
	return http.ListenAndServe(listenAllInterfaces, nil)
}

// return address of the proxy
func getProxyServer(proxyConfig *configs.Config) (string, error) {
	if (*proxyConfig).Proxy.Listen.Address == "" || (*proxyConfig).Proxy.Listen.Port == "" {
		return "", errors.New("Reverse proxy address or port not valid")
	}
	return (*proxyConfig).Proxy.Listen.Address + ":" + (*proxyConfig).Proxy.Listen.Port, nil
}

// main function
func Execute() {
	// Parse command line args to get path of config and load belancing strategy
	kingpin.New(filepath.Base(os.Args[0]), "A simple reverse proxy")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Read config YAML file
	log.Println("Loading configuration in", *configFile)
	proxyConfig, err := configs.Load(*configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing config file %s", *configFile))
		os.Exit(2)
	}

	// Start proxy
	if err := startProxy(&proxyConfig); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "Error starting proxy"))
		os.Exit(2)
	}
}
