package prepare

import "go.uber.org/zap"

func ZapLogger() *zap.SugaredLogger {
	logger := zap.NewExample().Sugar()

	defer func(sugar *zap.SugaredLogger) {
		_ = sugar.Sync()
	}(logger)

	return logger
}
