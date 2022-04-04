# amqexample   

## Introduction

This test code is for use as part of a blog post to be published on the xphyr.net web site. The code can be used to create an AMQP queue that fills faster than it is emptied, so that kubernetes autoscaling tools can be used to scale up the consumer in response to a backpressure of work to do.

## Building the code

```
$ git clone 
$ cd amqexample
$ go build
```

## Building a Container

```
$ podman build -t amqexample:latest .
```
# amqpexample
