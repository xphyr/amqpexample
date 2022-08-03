# amqexample

## Introduction

This test code is for use as part of a blog post to be published on the xphyr.net web site. The code can be used to create an AMQP queue that fills faster than it is emptied, so that kubernetes autoscaling tools can be used to scale up the consumer in response to a back pressure of work to do. The application is both the "sender" and the "receiver" in order to ease the deployment of this example application. This code is not intended to be an example of how to make a good AMQP application, just to get the job done.

This application has been tested against [ActiveMQ](https://activemq.apache.org/), [ActiveMQ Artemis](https://activemq.apache.org/components/artemis/) and [RabbitMQ](https://www.rabbitmq.com/).

## Building the code

The code can be built locally using the `go build` command. You will need to have the Go compiler installed in order to build it locally. If you do not have the Go compiler installed on your machine, you can also use the [container building process](#building-a-container) below.

```bash
$ git clone https://github.com/xphyr/amqpexample.git
$ cd amqpexample
$ go build
```

## Building a Container

The example application can also be built as a container.

> **NOTE:** The example command below uses podman, but you can substitute the docker command for podman with no issues.

```
$ podman build -t amqexample:latest .
# If you want to use your own image you will need to push the image to a registry
$ podman push amqexample:latest <container registry>/amqexample:latest
```

## Deploying the Application

The single Image can be used as both the application that creates and fills the AMQP queue and the application that removes the entries from the queue. Sample YAML is included in the `k8s` directory to deploy a pre-compiled application from the Quay registry. 