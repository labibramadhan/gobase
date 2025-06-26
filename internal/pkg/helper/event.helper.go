package helper

import (
	"context"
	"encoding/json"
	"strings"

	pkgErr "clodeo.tech/public/go-universe/pkg/err"
	"clodeo.tech/public/go-universe/pkg/tracer"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog/log"
)

/*
WrapProcessMessages processes messages from a channel, handles payloads, and logs the progress.
Parameters:
  - messages: Channel from which messages are received.
  - handlerFunc: Function to handle the unmarshaled payload.
  - spanName: Name of the tracing span for observability.
  - note: see the config router wtermill for now retry count if msg error and this is auto ack if use router
*/
func WrapProcessMessages[T any](msg *message.Message, handlerFunc func(ctx context.Context, payload T) error, spanName string) error {
	// Create a new tracing span for observability
	ctx, span := tracer.StartSpan(context.Background(), spanName, nil)
	defer span.End()
	// Process each message from the channel
	return processMessage(ctx, msg, handlerFunc, spanName)
}

/*
processMessage processes a single message, error handling, and logging.
Parameters:
  - ctx: Context for tracing and logging.
  - msg: The message to be processed.
  - handlerFunc: Function to handle the unmarshaled payload.
  - spanName: Name of the tracing span for observability.
*/
func processMessage[T any](ctx context.Context, msg *message.Message, handlerFunc func(ctx context.Context, payload T) error, spanName string) error {
	// Unmarshal the message payload
	var payload T
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		logErrorPayload(msg.Payload, err)
		return err
	}

	// Log the message handling start
	logMessage("PROCESS HANDLE MESSAGE", payload, spanName)

	// Handle the unmarshaled payload
	if err := handlerFunc(ctx, payload); err != nil {
		logErrorHanldeFunc(err)
		// ignored error for no retry msg
		cusErr := pkgErr.GetError(err)
		if cusErr.Type == pkgErr.ErrBadRequest ||
			strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "sql: no rows in result set") ||
			strings.Contains(err.Error(), "tag validation failed") {
			return nil
		}
		return err
	}

	// Log the message handling completion
	logMessage("DONE HANDLE MESSAGE", payload, spanName)
	return nil
}

/*
logErrorPayload logs an error encountered during payload unmarshaling.
Parameters:
  - payload: The raw payload data.
  - err: The error encountered during unmarshaling.
*/
func logErrorPayload(payload []byte, err error) {
	log.Error().Interface("[ERROR PAYLOAD]", string(payload)).Msgf("Error unmarshaling payload: %s", err.Error())
}

/*
logErrorHanldeMessage logs an error encountered during message handling.
Parameters:
  - err: The error encountered during message handling.
*/
func logErrorHanldeFunc(err error) {
	log.Error().Msgf("[ERROR HANDLE MESSAGE]: %s", err.Error())
}

/*
logMessage logs a message with the provided stage and span name for observability.
Parameters:
  - stage: The current stage of message processing.
  - payload: The unmarshaled payload data.
  - spanName: The name of the tracing span.
*/
func logMessage[T any](stage string, payload T, spanName string) {
	log.Info().Interface("["+stage+"]", payload).Msg(spanName)
}
