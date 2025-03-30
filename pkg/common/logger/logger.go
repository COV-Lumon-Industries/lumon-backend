package logger

import (
	"go.uber.org/zap"
)

var APILogger *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	APILogger = logger.Sugar()
}
