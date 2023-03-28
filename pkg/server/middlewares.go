package server

import (
	"colorblinder/pkg/metrics"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (s *Server) LoggingMiddleware(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()
		err := f(c)
		if c.Path() == "/stream*" && err == echo.ErrNotFound {
			return err
		}
		if err != nil {
			s.l.Warn("error processing request", zap.String("path", c.Path()),
				zap.Int("status", c.Response().Status), zap.Error(err))
		} else {
			s.l.Info("processed request", zap.String("path", c.Path()),
				zap.Duration("duration", time.Since(startTime)))
		}
		metrics.ReqTimeHist.WithLabelValues(strconv.Itoa(c.Response().Status)).Observe(time.Since(startTime).Seconds())
		return err
	}
}

func (s *Server) ActiveFiltersMiddleware(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() != "/stream*" {
			return f(c)
		}
		filterID := strings.Split(c.ParamValues()[0], "/")[1]
		s.c.UpdateLastExecutionTime(filterID, time.Now())

		return f(c)
	}
}
