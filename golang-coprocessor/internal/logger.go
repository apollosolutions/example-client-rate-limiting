package coprocessor

import (
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var logger logr.Logger = defaultLogger()

func defaultLogger() logr.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.MessageFieldName = "message"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	w := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
	zl := zerolog.New(w).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	zerologr.VerbosityFieldName = ""
	return zerologr.New(&zl)
}
