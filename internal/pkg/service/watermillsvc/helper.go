package watermillsvc

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func BuildMessage(messageId uuid.UUID, payload interface{}) (*message.Message, error) {
	var payloadBytes []byte
	var err error
	if payload != nil {
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}
	return message.NewMessage(messageId.String(), payloadBytes), nil
}

func BuildNewMessage(payload interface{}) (*message.Message, error) {
	var payloadBytes []byte
	var err error
	if payload != nil {
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}
	return message.NewMessage(uuid.New().String(), payloadBytes), nil
}
