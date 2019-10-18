# reverseproxy

A simple HTTP reverse proxy and load balancer that in written in Golang.
Traefik integrates with your existing infrastructure components ([Docker](https://www.docker.com/), [Swarm mode](https://docs.docker.com/engine/swarm/), [Kubernetes](https://kubernetes.io), [Marathon](https://mesosphere.github.io/marathon/), [Consul](https://www.consul.io/), [Etcd](https://coreos.com/etcd/), [Rancher](https://rancher.com), [Amazon ECS](https://aws.amazon.com/ecs), ...) and configures itself automatically and dynamically.

---

:warning: Please be aware that the old configurations for Traefik v1.X are NOT compatible with the v2.X config as of now. If you're running v2, please ensure you are using a [v2 configuration](https://docs.traefik.io/).

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

## Maintainers

Binh Luong <tbinhluong@gmail.com>