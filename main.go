package main

import (
	"flag"
	"os"

	"github.com/jamiealquiza/envy"
	"github.com/streadway/amqp"

	"fmt"
	"time"
)

var (
	amqpURI    = flag.String("server", "amqp://guest:guest@pobox.xphyrlab.net", "server address and port")
	mode       = flag.String("mode", "publisher", "act as publisher or consumer")
	debugLevel = flag.Bool("debug", false, "enable debug messages")
)

func main() {

	envy.Parse("AMQP") // looks for AMQP_SERVER, AMQP_MODE, AMQP_DEBUG etc
	flag.Parse()

	connection, _ := amqp.Dial(*amqpURI)
	defer connection.Close()

	channel, _ := connection.Channel()
	defer channel.Close()
	durable, exclusive := false, false
	autoDelete, noWait := true, true
	q, _ := channel.QueueDeclare("test", durable, autoDelete, exclusive, noWait, nil)
	channel.QueueBind(q.Name, "#", "amq.topic", false, nil)

	switch *mode {
	case "publisher":
		go publish(channel, &q)
	case "consumer":
		go subscribe(channel, &q)
	default:
		fmt.Println("must specify mode to work")
		os.Exit(1)
	}

	select {}
}

func publish(channel *amqp.Channel, q *amqp.Queue) {
	timer := time.NewTicker(3 * time.Second)

	for t := range timer.C {
		msg := amqp.Publishing{
			DeliveryMode: 1,
			Timestamp:    t,
			ContentType:  "text/plain",
			Body:         []byte("Hello world"),
		}
		mandatory, immediate := false, false
		channel.Publish("amq.topic", "ping", mandatory, immediate, msg)
		fmt.Println("pushed data")
	}
}

func subscribe(channel *amqp.Channel, q *amqp.Queue) {
	autoAck, exclusive, noLocal, noWait := false, false, false, false
	messages, _ := channel.Consume(q.Name, "", autoAck, exclusive, noLocal, noWait, nil)
	multiAck := false
	for msg := range messages {
		fmt.Println("Body:", string(msg.Body), "Timestamp:", msg.Timestamp)
		fmt.Println("Doing  work for 10 seconds")
		time.Sleep(10 * time.Second)
		msg.Ack(multiAck)
	}
}
