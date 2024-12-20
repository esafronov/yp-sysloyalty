package routine

import (
	"context"
	"sync"

	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/esafronov/yp-sysloyalty/internal/logger"
	"github.com/esafronov/yp-sysloyalty/internal/usecase"
	"go.uber.org/zap"
)

type Updater struct {
	cr domain.CustomerRepository
	or domain.OrderRepository
}

// updater factory
func NewUpdater(or domain.OrderRepository, cr domain.CustomerRepository, params *config.AppParams) *Updater {
	return &Updater{
		cr: cr,
		or: or,
	}
}

// runs rupdater routine
func (u *Updater) Run(ctx context.Context, updateChan chan *domain.OrderUpdate, wg *sync.WaitGroup) {
	uc := usecase.NewOrdersUpdateUsecase(u.or, u.cr)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for update := range updateChan {
			logger.Log.Info("received update",
				zap.String("order", update.Num),
				zap.String("status", string(update.Status)),
				zap.Int64("accrual", update.Accrual),
			)
			if err := uc.Update(ctx, update); err != nil {
				logger.Log.Error("update err", zap.String("order", update.Num), zap.Error(err))
				continue
			}
		}
		logger.Log.Info("exit from updater...")
	}()
}
