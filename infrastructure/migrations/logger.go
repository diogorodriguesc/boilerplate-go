package migrations

import "github.com/rs/zerolog/log"

type CustomLogger struct{}

func (l CustomLogger) Print(v ...interface{}) {
	log.Info().Msgf("%v", v...)
}

func (l CustomLogger) Println(v ...interface{}) {
	log.Info().Msgf("%v", v...)
}

func (l CustomLogger) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func (l CustomLogger) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}
