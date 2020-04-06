# Kafka Web App

This is a little interactive web app to demonstrate Kafka-based real-time event streaming.

The `/produce` endpoint allows users to pick different Confluent logos, generating a stream of events.

The `/consume` endpoint displays a logo that changes in real-time based on the produced events.

This demo is web-scale, and admits multiple simultaneous producers and consumers.

## Usage

### Running Locally

You will need to install `librdkafka` and `go` first if you don't already have those on your
workstation. You will likely need to set or extend the `PKG_CONFIG_PATH` environment variable to
include a path to your directory containing `libcrypto.pc` for `librdkafka` to work. This package
comes with OpenSSL; Macs come with LibreSSL pre-installed which is insufficient. You may need to
`brew install openssl@1.1` where you will then see output like:

```
For pkg-config to find openssl@1.1 you may need to set:
  export PKG_CONFIG_PATH="/usr/local/opt/openssl@1.1/lib/pkgconfig"
```

Follow that instruction.

Next, start Kafka locally, for example by following the [Apache Kafka Quickstart](https://kafka.apache.org/quickstart).

Then run `./start-server.sh` from the root of this repository.

Finally, open up the [Consume](http://localhost:8080/consume) and [Produce](http://localhost:8080/produce)
pages in a couple browser windows, and follow the instructions on the pages.

## TODO

1. Add instructions for running this app and Kafka on Minikube
1. Add instructions for running on Minikube and connecting it to a Confluent Cloud (CC) cluster
1. Publish Docker images
1. Add instructions for running this app and Kafka on a "real" Kubernetes cluster
1. Add instructions for running on a "real" Kubernetes cluster and connecting it to a CC cluster
1. Code TODOs
