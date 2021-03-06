package main

import (
	"github.com/spf13/viper"
	"github.com/streadway/amqp"

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
func (ms MessageService) SendInitiation(i Initiation) (err error) {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var data []byte

	if data, err = json.Marshal(i); err != nil {
		return err
	}

	if conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/"); err != nil {
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

	if err = ch.Publish(
		"",
		initiationQueue,
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
