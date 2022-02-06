package agent

import (
	"os"

	"github.com/rs/zerolog"
)

func NewAgentLogger() (logger *zerolog.Logger, err error) {

	logFile, err := os.OpenFile("agent_log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		return nil, err
	}
	lg := zerolog.New(logFile).With().Timestamp().Logger()
	logger = &lg
	return logger, nil
}
