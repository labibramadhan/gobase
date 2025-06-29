package watermillsvc

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog"
)

// NewZerologLoggerAdapter creates a Watermill logger adapter for zerolog.
func NewZerologLoggerAdapter(logger zerolog.Logger) watermill.LoggerAdapter {
	return &zerologLoggerAdapter{logger: logger}
}

type zerologLoggerAdapter struct {
	logger zerolog.Logger
}

func (l *zerologLoggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	l.logger.Error().Err(err).Fields(fields).Msg(msg)
}

func (l *zerologLoggerAdapter) Info(msg string, fields watermill.LogFields) {
	l.logger.Info().Fields(fields).Msg(msg)
}

func (l *zerologLoggerAdapter) Debug(msg string, fields watermill.LogFields) {
	l.logger.Debug().Fields(fields).Msg(msg)
}

func (l *zerologLoggerAdapter) Trace(msg string, fields watermill.LogFields) {
	l.logger.Trace().Fields(fields).Msg(msg)
}

func (l *zerologLoggerAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &zerologLoggerAdapter{logger: l.logger.With().Fields(fields).Logger()}
}
