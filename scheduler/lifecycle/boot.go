package lifecycle

import (
	"context"
	"go.uber.org/zap"
	"scheduler/internal/logger"
)

type Bootable interface {
	Boot(context.Context) error
}

func BootAll(ctx context.Context, components ...Bootable) error {
	for _, component := range components {
		if err := component.Boot(ctx); err != nil {
			logger.Log.Error("Boot Error", zap.Error(err))
		}
	}
	return nil
}
