package main

import (
	"colorblinder/pkg/config"
	"colorblinder/pkg/logger"
	"colorblinder/pkg/metrics"
	"colorblinder/pkg/server"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.Provide(
			logger.ProvideLogger,
			config.ProvideConfig,
		),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			server.ProvideServer,
		),
		fx.Invoke(
			server.InvokeServer,
			metrics.InvokeMetricsServer,
		),
	)
	app.Run()
}
