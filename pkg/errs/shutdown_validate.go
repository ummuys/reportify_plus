package errs

import "github.com/rs/zerolog"

type SDMsg struct {
	Err  error
	From string
}

func ShutdownStatus(logger zerolog.Logger, SDChan <-chan SDMsg) {
	hadErr := false
	for msg := range SDChan {
		logger.Error().Err(msg.Err).Msg(msg.From)
		hadErr = true
	}

	if hadErr {
		logger.Fatal().Msg("shutdown with errors")
	} else {
		logger.Info().Msg("graceful shutdown")
	}
}
