# Apache NiFi Exporter

![go report](https://goreportcard.com/badge/github.com/diarworld/nifi_exporter)
![version](https://img.shields.io/docker/v/diarworld/nifi_exporter?sort=semver)
![image size](https://img.shields.io/docker/image-size/diarworld/nifi_exporter?sort=semver)

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

You may enable or disable some collectors using config file.

See the [sample config](./sample-config.yml) for a full example of all available options.

## Running

### Using Docker

Docker image is available at [Docker Hub](https://hub.docker.com/diarworld/nifi_exporter):

```sh
docker run -p 9103:9103 -v /path/to/config.yml:/config/config.yml:ro diarworld/nifi_exporter:0.3.0
```

### Without Docker

Download a release package for your system from [Releases page](https://github.com/diarworld/nifi_exporter/releases), unpack it and run the binary directly:

```sh
curl -fLO https://github.com/diarworld/nifi_exporter/releases/download/v0.3.0/nifi_exporter-0.3.0.linux-amd64.tar.gz
tar -xvf nifi_exporter-0.3.0.linux-amd64.tar.gz
cd ./nifi_exporter-0.3.0.linux-amd64
./nifi_exporter /path/to/config.yml
```

## Building

```sh
go get github.com/diarworld/nifi_exporter
cd ${GOPATH-$HOME/go}/src/github.com/diarworld/nifi_exporter
go build
./nifi_exporter
```
