package services

import (
	"github.com/alanwade2001/spa-common/rabbitmq"
	"github.com/alanwade2001/spa-submissions-api/models/generated/initiation"
	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/spf13/viper"

	"encoding/json"
)

// MessageService s
type MessageService struct {
	messaging *rabbitmq.Messaging
}

// NewMessageService f
func NewMessageService() types.MessageAPI {
	return MessageService{}
}

func (ms *MessageService) GetMessaging() *rabbitmq.Messaging {
	if ms.messaging != nil {
		return ms.messaging
	}

	messageServiceURI := viper.GetString("MESSAGE_SERVICE_URI")
	initiationQueue := viper.GetString("INITIATION_QUEUE")

	ms.messaging = &rabbitmq.Messaging{
		Url:       messageServiceURI,
		QueueName: initiationQueue,
	}

	return ms.messaging
}

// SendInitiation f
func (ms MessageService) SendInitiation(i initiation.InitiationModel) (err error) {
	var data []byte

	if data, err = json.Marshal(i); err != nil {
		return err
	}

	if err := ms.GetMessaging().Connect(); err != nil {
		return nil
	}

	defer ms.GetMessaging().Disconnect()

	if err := ms.GetMessaging().Publish("application/json", data); err != nil {
		return nil
	}

	return nil
}
