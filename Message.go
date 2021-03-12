package main

import (
	types "github.com/alanwade2001/spa-common"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"k8s.io/klog/v2"

	"encoding/json"
)

// MessageService s
type MessageService struct {
}

// NewMessageService f
func NewMessageService() MessageAPI {
	return MessageService{}
}

// SendInitiation f
func (ms MessageService) SendInitiation(i types.Initiation) (err error) {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var q amqp.Queue
	var data []byte

	if data, err = json.Marshal(i); err != nil {
		return err
	}

	messageServiceURI := viper.GetString("MESSAGE_SERVICE_URI")
	if conn, err = amqp.Dial(messageServiceURI); err != nil {
		return err
	}
	defer conn.Close()

	// Let's start by opening a channel to our RabbitMQ instance
	// over the connection we have already established
	if ch, err = conn.Channel(); err != nil {
		return err
	}
	defer ch.Close()

	initiationQueue := viper.GetString("INITIATION_QUEUE")
	klog.Infof("initiation queue:[%s]", initiationQueue)

	if q, err = ch.QueueDeclare(
		initiationQueue, // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	); err != nil {
		return err
	}

	if err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	); err != nil {
		return err
	}

	return nil
}
