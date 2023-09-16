package logmiddleware

import "go.uber.org/zap"

func Connect(production bool) *zap.Logger {

	var logger *zap.Logger

	if production == true {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}

	return logger
}
