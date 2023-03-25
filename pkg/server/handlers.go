package server

import (
	"colorblinder/internal/filter"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CreateFilterRequest struct {
	RGBAOverlay      [4]int `json:"rgba_overlay,omitempty"`
	StartSecond      int    `json:"start_second,omitempty"`
	IsPhotosensitive bool   `json:"is_photosensitive,omitempty"`
}

type CreateFilterResponse struct {
	ID string `json:"id"`
}

func (s *Server) CreateFilter(c echo.Context) error {
	var req CreateFilterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	f := filter.NewFilter(req.RGBAOverlay, req.StartSecond, req.IsPhotosensitive)
	s.filters[f.ID] = f
	return c.JSON(http.StatusOK, CreateFilterResponse{ID: f.ID})
}

type StartStreamRequest struct {
	FilterID  string `json:"filter_id"`
	StreamURL string `json:"stream_url"`
}

type StartStreamResponse struct {
	NewURL string `json:"new_url"`
}

func (s *Server) StartStream(c echo.Context) error {
	var req StartStreamRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	f, ok := s.filters[req.FilterID]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "filter not found")
	}
	ctx := context.Background()
	go func(url string) {
		err := f.StartProcess(url)
		if err != nil {
			s.l.Warn("error processing", zap.Error(err))
		}
		ctx.Done()
	}(req.StreamURL)
	go func() {
		<-ctx.Done()
		s.l.Info("killed stream process", zap.String("id", f.ID))
		delete(s.filters, f.ID)
	}()
	return c.JSON(http.StatusOK, StartStreamResponse{NewURL: "/stream/" + req.FilterID + "/file.mpd"})
}
