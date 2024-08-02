package logger

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

var (
	Logs zerolog.Logger
	file *os.File
	once sync.Once
)

func initLogger() {
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

// GetLogger returns the singleton instance of the logger
func init() {
	once.Do(initLogger)
	err := godotenv.Load(".env")
	if err != nil {
		Logs.Fatal().Msg("Error loading .env file")
	}
}
