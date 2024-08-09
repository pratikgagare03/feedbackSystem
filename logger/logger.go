package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var (
	Logs zerolog.Logger
	file *os.File
)

func StartLogger() {
	var err error
	file, err = os.OpenFile(
		"myapp.log",
		os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	Logs = zerolog.New(file).With().Timestamp().Logger()
}
