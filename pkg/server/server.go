package server

import (
	"colorblinder/internal/filter"
	"colorblinder/pkg/config"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	Address string
	e       *echo.Echo
	l       *zap.Logger
	filters map[string]filter.Filter
}

func ProvideServer(config *config.Config, l *zap.Logger) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Debug = false
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	e.Use()
	e.Static("/stream", "/tmp")
	return &Server{
		Address: config.ServerConfig.Address,
		e:       e,
		l:       l,
		filters: make(map[string]filter.Filter),
	}
}

func (s *Server) Start() error {
	s.e.POST("/create_filter", s.CreateFilter)
	s.e.POST("/start_stream", s.StartStream)

	s.e.Use(s.LoggingMiddleware)
	return s.e.Start(s.Address)
}

func InvokeServer(lifecycle fx.Lifecycle, s *Server) {
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if e := s.Start(); e != nil {
					panic(e)
				}
			}()
			return nil
		},
	})
}
