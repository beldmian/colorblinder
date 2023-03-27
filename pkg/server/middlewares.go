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
			activeFiltersIDs := []string{}
			for _, filter := range s.activeFilters {
				activeFiltersIDs = append(activeFiltersIDs, filter.ID)
			}
			s.l.Info("processed request", zap.String("path", c.Path()),
				zap.Duration("duration", time.Since(startTime)),
				zap.Strings("active_filters", activeFiltersIDs))
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
		fragment_path := c.ParamValues()[0]
		time_now := time.Now()
		for id, filter := range s.activeFilters {
			if strings.Contains(fragment_path, filter.ID) {
				s.activeFiltersMu.Lock()
				filter.LastExecutionTime = time_now
				s.activeFilters[id] = filter
				s.activeFiltersMu.Unlock()
				continue
			}
			if time.Since(filter.LastExecutionTime) > time.Second*30 {
				s.activeFiltersMu.Lock()
				filter.ContextCancel()
				delete(s.activeFilters, id)
				s.activeFiltersMu.Unlock()
				s.l.Info("killed filter due to inactivity", zap.String("filter_id", filter.ID))
			}
		}
		return f(c)
	}
}
