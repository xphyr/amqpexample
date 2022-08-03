package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/Azure/go-amqp"
	"github.com/jamiealquiza/envy"
)

var (
	amqpURI    = flag.String("server", "amqp://guest:guest@pobox.xphyrlab.net", "server address and port")
	mode       = flag.String("mode", "publisher", "act as publisher or consumer")
	debugLevel = flag.Bool("debug", false, "enable debug messages")
)

func main() {

	envy.Parse("AMQP") // looks for AMQP_SERVER, AMQP_MODE, AMQP_DEBUG etc
	flag.Parse()

	// Create client
	client, err := amqp.Dial(*amqpURI)
	if err != nil {
		log.Fatal("Dialing AMQP server:", err)
	}
	defer client.Close()

	// Open a session
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Creating AMQP session:", err)
	}

	switch *mode {
	case "publisher":
		go publish(session)
	case "consumer":
		go subscribe(session)
	default:
		fmt.Println("must specify mode to work")
		os.Exit(1)
	}

	select {}
}

func publish(session *amqp.Session) {

	ctx := context.Background()
	i := 0

	// Create a sender
	sender, err := session.NewSender(
		amqp.LinkTargetAddress("keda-test"),
	)
	if err != nil {
		log.Fatal("Creating sender link:", err)
	}

	timer := time.NewTicker(3 * time.Second)

	for t := range timer.C {
		myMessage := "Hello " + t.String()
		err = sender.Send(ctx, amqp.NewMessage([]byte(myMessage)))
		if err != nil {
			log.Fatal("Sending message:", err)
		}
		i = i + 1
	}
	sender.Close(ctx)
}

func subscribe(session *amqp.Session) {
	ctx := context.Background()
	// Create a receiver
	receiver, err := session.NewReceiver(
		amqp.LinkSourceAddress("keda-test"),
		amqp.LinkCredit(10),
	)
	if err != nil {
		log.Fatal("Creating receiver link:", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		receiver.Close(ctx)
		cancel()
	}()

	for {
		// Receive next message
		msg, err := receiver.Receive(ctx)
		if err != nil {
			log.Fatal("Reading message from AMQP:", err)
		}

		// Accept message
		receiver.AcceptMessage(ctx, msg)

		fmt.Printf("Message received: %s\n", msg.GetData())
		time.Sleep(5 * time.Second)
	}
}
