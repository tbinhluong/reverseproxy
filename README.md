<h1 align="center">
    reverseproxy
</h1>

<p align="center">
    A simple HTTP reverse proxy and load balancer that in written in Golang.
    <br/><br/>
    <a href="https://github.com/tbinhluong/reverseproxy/releases">
        <img alt="latest version" src="https://img.shields.io/github/tag/tbinhluong/reverseproxy.svg" />
    </a>
    <a href="https://www.apache.org/licenses/LICENSE-2.0">
        <img alt="Apache-2.0 License" src="https://img.shields.io/github/license/tbinhluong/reverseproxy.svg" />
    </a>
    <a href="https://microbadger.com/images/tbinhluong/reverseproxy">
      <img alt="Microbadger" src="https://images.microbadger.com/badges/image/tbinhluong/reverseproxy.svg" />
    </a>
    <a href="https://hub.docker.com/r/tbinhluong/reverseproxy">
        <img alt="Pulls from DockerHub" src="https://img.shields.io/docker/pulls/tbinhluong/reverseproxy.svg?style=flat-square" />
    </a>
</p>
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
docker run -d -p 8080:8080  -v $PWD/config.yml:/reverseproxy/config.yml tbinhluong/reverseproxy:latest
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