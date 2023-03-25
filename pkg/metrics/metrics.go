package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

var (
	ReqTimeHist = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "req_time_hist",
		Buckets: []float64{1, 2, 5, 6, 10},
	}, []string{"code"})
	ActiveFilterers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_filterers",
	})
)

func InvokeMetricsServer(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				http.Handle("/metrics", promhttp.Handler())
				_ = http.ListenAndServe(":2112", nil)
			}()
			return nil
		},
	})
}
