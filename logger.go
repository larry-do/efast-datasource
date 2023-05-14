package godatasource

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

var (
	loggingWriter logger.Writer
)

func LoggingWriter(writer logger.Writer) {
	loggingWriter = writer
}

type ZerologWrappingWriter struct {
	logger.Writer
	Logger zerolog.Logger
}

func (writer ZerologWrappingWriter) Printf(msg string, data ...interface{}) {
	writer.Logger.Info().Msgf(msg, data...)
}
