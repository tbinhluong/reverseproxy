# reverseproxy

A simple HTTP reverse proxy and load balancer that in written in Golang.
---

## Features

- Listens to HTTP requests and forward them to downstream services
- Supports 2 load balancing algorithms: randomly and round-robin
- Provides endpoint /healthz for health check
- Provides endpoint /metrics to indicate the availability SLI of the reverse proxy
- Packaged as a single binary file (made with with go) and available as a [docker image](https://hub.docker.com/r/tbinhluong/reverseproxy)
- Packaged as a [helm chart](https://github.com/tbinhluong/tbinhluong.github.io/tree/master/charts/reverseproxy-helm) 


## Getting Started

- Grab the latest binary from the [releases](https://github.com/tbinhluong/reverseproxy/releases) page and run it with the [sample configuration file](https://raw.githubusercontent.com/tbinhluong/reverseproxy/master/config/config.yml):

```shell
./reverseproxy --config.file=config.yml
```

- Or use the official tiny Docker image and run it with the [sample configuration file](https://raw.githubusercontent.com/tbinhluong/reverseproxy/master/config/config.yml):

```shell
docker run -d -p 8080:8080  -v $PWD/config.yml:/reverseproxy/config.yml reverseproxy
```

- Or use the helm chart

```shell
helm repo add myrepo https://tbinhluong.github.io/
helm install myrepo/reverseproxy-helm --name=reverseproxy
```

- Or get the sources:

```shell
git clone https://github.com/tbinhluong/reverseproxy
```

## Maintainers

Binh Luong <tbinhluong@gmail.com>