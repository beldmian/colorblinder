package server

import (
	"colorblinder/internal/filter"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type StartStreamRequest struct {
	StreamURL string `json:"stream_url"`
}

type StartStreamResponse struct {
	NewURL string `json:"new_url"`
}

func (s *Server) StartStream(c echo.Context) error {
	var req StartStreamRequest
	pid := uuid.NewString()
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func(url string) {
		err := filter.StartProcess(ctx, pid, url)
		if err != nil {
			s.l.Warn("error processing", zap.Error(err))
		}
	}(req.StreamURL)
	go func() {
		<-ctx.Done()
		s.l.Info("killed stream process", zap.String("id", pid))
	}()
	s.activeFilters[pid] = FilterInfo{
		ID:                pid,
		ContextCancel:     cancel,
		LastExecutionTime: time.Now(),
	}
	return c.JSON(http.StatusOK, StartStreamResponse{NewURL: "/stream/" + pid + "/file.mpd"})
}
