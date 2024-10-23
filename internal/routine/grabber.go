package routine

import (
	"context"
	"sync"
	"time"

	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/esafronov/yp-sysloyalty/internal/logger"
	"go.uber.org/zap"
)

type Grabber struct {
	or           domain.OrderRepository
	grabInterval time.Duration
}

func NewGrabber(or domain.OrderRepository, params *config.AppParams) *Grabber {
	return &Grabber{
		or:           or,
		grabInterval: time.Duration(1) * time.Minute,
	}
}

func (g *Grabber) Run(appCtx context.Context, retryAfterChan <-chan int, wg *sync.WaitGroup) chan *domain.Order {
	ch := make(chan *domain.Order, 1)
	ticker := time.NewTicker(g.grabInterval)
	var delayGrabber = func(pause int) {
		logger.Log.Info("delay collecting orders for a while...", zap.Int("pause", pause))
		ticker.Reset(time.Duration(pause) * time.Minute)
	}
	wg.Add(1)
	go func() {
		defer func() {
			logger.Log.Info("exit from grabber...")
			ticker.Stop()
			close(ch)
			wg.Done()
		}()
		var retryAfter int
	out:
		for {
			select {
			case <-appCtx.Done():
				break out
			case retryAfter = <-retryAfterChan:
				delayGrabber(retryAfter)
			case <-ticker.C:
				if retryAfter > 0 {
					ticker.Reset(g.grabInterval)
				}
				logger.Log.Info("collect orders for update...")
				orders, err := g.or.GetNotFinished(appCtx)
				if err != nil {
					logger.Log.Error("collect orders for update", zap.Error(err))
					continue
				}
				for _, order := range orders {
					select {
					case <-appCtx.Done():
						break out
					case retryAfter = <-retryAfterChan:
						delayGrabber(retryAfter)
						return
					default:
						ch <- order
					}
				}
			}
		}
	}()
	return ch
}
