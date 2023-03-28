package cleaner

import (
	"context"
	"sync"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type FilterInfo struct {
	ID                string
	ContextCancel     context.CancelFunc
	LastExecutionTime time.Time
}

type Cleaner struct {
	l                 *zap.Logger
	activeFiltersInfo map[string]FilterInfo
	activeFiltersMu   sync.Mutex
}

func ProvideCleaner(l *zap.Logger) *Cleaner {
	return &Cleaner{
		l:                 l,
		activeFiltersInfo: make(map[string]FilterInfo, 0),
		activeFiltersMu:   sync.Mutex{},
	}
}

func (c *Cleaner) Start() {
	for {
		for _, filter := range c.activeFiltersInfo {
			if time.Since(filter.LastExecutionTime) > time.Second*30 {
				c.l.Info("cleaning filter", zap.String("id", filter.ID))
				filter.ContextCancel()
				c.activeFiltersMu.Lock()
				delete(c.activeFiltersInfo, filter.ID)
				c.activeFiltersMu.Unlock()
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func InvokeCleaner(lifecycle fx.Lifecycle, c *Cleaner) {
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go c.Start()
			return nil
		},
	})
}

func (c *Cleaner) UpdateLastExecutionTime(id string, lastExecutionTime time.Time) {
	if filter, ok := c.activeFiltersInfo[id]; ok {
		filter.LastExecutionTime = lastExecutionTime
		c.activeFiltersMu.Lock()
		c.activeFiltersInfo[id] = filter
		c.activeFiltersMu.Unlock()
	}
}

func (c *Cleaner) AddFilter(filter FilterInfo) {
	c.activeFiltersMu.Lock()
	c.activeFiltersInfo[filter.ID] = filter
	c.activeFiltersMu.Unlock()
}
