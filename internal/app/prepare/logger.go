package prepare

import "go.uber.org/zap"

func ZapLogger() *zap.SugaredLogger {
	zapLogger, _ := zap.NewProduction()
	logger := zapLogger.Sugar()

	defer func(sugar *zap.SugaredLogger) {
		_ = sugar.Sync()
	}(logger)

	return logger
}
