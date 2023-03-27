package server

import (
	"colorblinder/pkg/config"
	"context"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type FilterInfo struct {
	ID                string
	ContextCancel     context.CancelFunc
	LastExecutionTime time.Time
}

type Server struct {
	Address         string
	e               *echo.Echo
	l               *zap.Logger
	activeFilters   map[string]FilterInfo
	activeFiltersMu sync.Mutex
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
	e.Static("/stream", "/tmp")
	return &Server{
		Address:         config.ServerConfig.Address,
		e:               e,
		l:               l,
		activeFilters:   make(map[string]FilterInfo),
		activeFiltersMu: sync.Mutex{},
	}
}

func (s *Server) Start() error {
	s.e.POST("/start_stream", s.StartStream)

	s.e.Use(s.LoggingMiddleware)
	s.e.Use(s.ActiveFiltersMiddleware)
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
