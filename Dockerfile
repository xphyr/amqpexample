##
## Build
##
FROM docker.io/golang as amqpexample_builder
WORKDIR /opt/app-root
ADD . /opt/app-root
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/amqpexample

##
## Deploy
##
FROM docker.io/alpine:latest
COPY --from=amqpexample_builder /opt/app-root/bin/amqpexample /usr/bin/amqpexample
ENTRYPOINT [ "amqpexample" ]