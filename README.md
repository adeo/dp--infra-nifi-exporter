# Apache NiFi Exporter

![go report](https://goreportcard.com/badge/github.com/adeo/dp--infra-nifi-exporter)
![version](https://img.shields.io/docker/v/adeo/dp--infra-nifi-exporter?sort=semver)
![image size](https://img.shields.io/docker/image-size/adeo/dp--infra-nifi-exporter?sort=semver)

Exports metrics from Apache NiFi API in Prometheus-compatible format.

## Configuration

NiFi Exporter is configured through a single YAML file. Here is the minimal configuration:

```yaml
---
exporter:
  listenAddress: 127.0.0.1:9103
nodes:
  - url: http://localhost:8080
    username: xxxxxx
    password: xxxxxx
```

See the [sample config](./sample-config.yml) for a full example of all available options.

## Running

### Using Docker

Docker image is available at [Docker Hub](https://hub.docker.com/adeo/dp--infra-nifi-exporter):

```sh
docker run -p 9103:9103 -v /path/to/config.yml:/config/config.yml:ro adeo/dp--infra-nifi-exporter:0.2.0
```

### Without Docker

Download a release package for your system from [Releases page](https://github.com/adeo/dp--infra-nifi-exporter/releases), unpack it and run the binary directly:

```sh
curl -fLO https://github.com/adeo/dp--infra-nifi-exporter/releases/download/v0.2.0/nifi_exporter-0.2.0.linux-amd64.tar.gz
tar -xvf nifi_exporter-0.2.0.linux-amd64.tar.gz
cd ./nifi_exporter-0.2.0.linux-amd64
./nifi_exporter /path/to/config.yml
```

## Building

```sh
go get github.com/adeo/dp--infra-nifi-exporter
cd ${GOPATH-$HOME/go}/src/github.com/adeo/dp--infra-nifi-exporter
go build
./nifi_exporter
```
