package main

import (
	"fmt"
	"log"

	"github.com/go-retail/pos-server/pkg/utils"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var rmq RMQ

//RMQ ..
type RMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      *amqp.Queue
}

func initRMQ() {
	//TODO externalize RabbitMQ Port to config
	urlString := fmt.Sprintf("amqp://%s:%s@%s:5672", viper.GetString("username"), viper.GetString("password"), viper.GetString("host"))
	log.Printf(urlString)
	conn, err := amqp.Dial(urlString)
	utils.FailOnError(err, "Failed to Connect to RabbitMQ")

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to Open The Channel")

	//TODO Deliver to Exchange not a Queue
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	utils.FailOnError(err, "Unable to Declare a Queue")
	rmq = RMQ{conn, ch, &q}
}

// Find a place to put this code ...

// defer rmq.Connection.Close()
// defer rmq.Channel.Close()
