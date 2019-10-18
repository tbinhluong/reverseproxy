// Copyright © 2019 Binh Luong
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
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
	configFile    = kingpin.Flag("config.file", "Path of configuration YAML file.").Default("config.yml").String()
	roundRobin    = kingpin.Flag("roundrobin", "Enable round-robin as load balancing strategy, otherwise randomly").Default("false").Bool()
	goodResponses uint32
	totalRequests uint32
)

type Metrics struct {
	GoodResponses uint32
	TotalRequests uint32
	SLI           string
}

type wrapperResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// A wrapper of http.ResponseWriter to get status code of response
func NewWrapperResponseWriter(w http.ResponseWriter) *wrapperResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &wrapperResponseWriter{w, http.StatusOK}
}

func (res *wrapperResponseWriter) WriteHeader(code int) {
	res.statusCode = code
	res.ResponseWriter.WriteHeader(code)
}

// handler sends requests to a service instance and forwards response back to origin
func handler(res http.ResponseWriter, req *http.Request) {

	// increase the number of total requests by 1 to measure SLI
	totalRequests++

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
		if getMetrics(req) {
			// goodResponses+1 because it is increased later
			sli := float32(goodResponses+1) / float32(totalRequests) * 100
			stringSLI := fmt.Sprintf("%f%%", sli)
			metrics := Metrics{GoodResponses: goodResponses + 1, TotalRequests: totalRequests, SLI: stringSLI}
			json.NewEncoder(res).Encode(metrics)
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

// Implement a new handler to get the status code of response
func wrapperHandler(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		lrw := NewWrapperResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, req)
		statusCode := lrw.statusCode
		log.Printf("<-- %d %s", statusCode, http.StatusText(statusCode))
		if statusCode < 500 {
			goodResponses++
		}
	})
}

func isHealthCheck(req *http.Request) bool {
	if req.RequestURI == "/healthz" && req.Method == "GET" {
		return true
	}
	return false
}

func getMetrics(req *http.Request) bool {
	if req.RequestURI == "/metrics" && req.Method == "GET" {
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

	wHandler := wrapperHandler(http.HandlerFunc(handler))
	http.Handle("/", wHandler)

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
