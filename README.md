# ambient-exporter

Prometheus exporter that collects device data from the Ambient Weather API.

![Build Docker image](https://github.com/ianunruh/ambient-exporter/actions/workflows/docker-build.yml/badge.svg)

Tested with the Ambient Weather WS-8482.

[Example metrics](docs/example-metrics.txt)

## Features

* Supports multiple base stations and multiple sensors per station
* Lightweight Docker image (around 15MB)

## Usage

Start by getting an API key and application key from [Ambient Weather](https://ambientweather.net/account).

```bash
export AMBIENT_API_KEY="abc"
export AMBIENT_APP_KEY="xyz"

docker compose up -d --build
docker compose ps
docker compose logs exporter

curl -s http://localhost:9090/metrics | grep ambient
```

## Deployment

Use Kustomize to deploy the exporter to a Kubernetes cluster.

```bash
kubectl -n monitoring create secret generic ambient-exporter \
    --from-literal=AMBIENT_API_KEY=${AMBIENT_API_KEY} \
    --from-literal=AMBIENT_APP_KEY=${AMBIENT_APP_KEY}

kubectl kustomize "https://github.com/ianunruh/ambient-exporter.git/deploy/basic?ref=v1.0.1" | \
    kubectl apply -n monitoring -f-
```
