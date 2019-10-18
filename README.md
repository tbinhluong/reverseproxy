# reverseproxy
---

A simple HTTP reverse proxy and load balancer that in written in Golang.
Traefik integrates with your existing infrastructure components ([Docker](https://www.docker.com/), [Swarm mode](https://docs.docker.com/engine/swarm/), [Kubernetes](https://kubernetes.io), [Marathon](https://mesosphere.github.io/marathon/), [Consul](https://www.consul.io/), [Etcd](https://coreos.com/etcd/), [Rancher](https://rancher.com), [Amazon ECS](https://aws.amazon.com/ecs), ...) and configures itself automatically and dynamically.

---

:warning: Please be aware that the old configurations for Traefik v1.X are NOT compatible with the v2.X config as of now. If you're running v2, please ensure you are using a [v2 configuration](https://docs.traefik.io/).

## Overview

Imagine that you have deployed a bunch of microservices with the help of an orchestrator (like Swarm or Kubernetes) or a service registry (like etcd or consul).
Now you want users to access these microservices, and you need a reverse proxy.

Traditional reverse-proxies require that you configure _each_ route that will connect paths and subdomains to _each_ microservice. 
In an environment where you add, remove, kill, upgrade, or scale your services _many_ times a day, the task of keeping the routes up to date becomes tedious. 

**This is when Traefik can help you!**

Traefik listens to your service registry/orchestrator API and instantly generates the routes so your microservices are connected to the outside world -- without further intervention from your part. 

**Run Traefik and let it do the work for you!** 
_(But if you'd rather configure some of your routes manually, Traefik supports that too!)_

![Architecture](docs/content/assets/img/traefik-architecture.png)

## Features

- Listens to HTTP requests and forward them to downstream services
- Supports 2 load balancing algorithms: randomly and round-robin
- Provides endpoint /healthz for health check
- Packaged as a single binary file (made with with go) and available as a [docker image](https://hub.docker.com/r/tbinhluong/reverseproxy)
- Packaged as a [helm chart](https://github.com/tbinhluong/tbinhluong.github.io/tree/master/charts/reverseproxy-helm) 


## Download

- Grab the latest binary from the [releases](https://github.com/tbinhluong/reverseproxy/releases) page and run it with the [sample configuration file](https://raw.githubusercontent.com/tbinhluong/reverseproxy/master/config/config.yml):

```shell
./reverseproxy --config.file=config.yml
```

- Or use the official tiny Docker image and run it with the [sample configuration file](https://raw.githubusercontent.com/tbinhluong/reverseproxy/master/config/config.yml):

```shell
docker run -d -p 8080:8080  -v $PWD/config.yml:/reverseproxy/config.yml reverseproxy
```

- Or get the sources:

```shell
git clone https://github.com/tbinhluong/reverseproxy
```

## Introductory Videos

You can find high level and deep dive videos on [videos.containo.us](https://videos.containo.us)

## Maintainers

Binh Luong <tbinhluong@gmail.com>