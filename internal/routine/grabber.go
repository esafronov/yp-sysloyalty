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
	grabInterval time.Duration //orders query interval
	selectLimit  int           //orders query limit
}

// grabber factory
func NewGrabber(or domain.OrderRepository, params *config.AppParams) *Grabber {
	return &Grabber{
		or:           or,
		grabInterval: time.Duration(*params.GrabInterval) * time.Second,
		selectLimit:  *params.GrabInterval * *params.ProcessRate,
	}
}

// runs grabber routine and returns order channel
func (g *Grabber) Run(appCtx context.Context, retryAfterChan <-chan int, wg *sync.WaitGroup) chan *domain.Order {
	ch := make(chan *domain.Order, 1)
	ticker := time.NewTicker(g.grabInterval)
	var delayGrabber = func(pause int) {
		logger.Log.Info("delay collecting orders for a while...", zap.Int("minutes", pause))
		ticker.Reset(time.Duration(pause*60) * time.Second)
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
			case <-appCtx.Done(): //app exit signal
				break out
			case retryAfter = <-retryAfterChan: //retryAfter signal from poll workers
				//delay grabbing process for retryAfter minutes
				delayGrabber(retryAfter)
			case <-ticker.C:
				if retryAfter > 0 {
					//restore normal ticker interval
					ticker.Reset(g.grabInterval)
				}
				logger.Log.Info("collect orders for update...")
				orders, err := g.or.GetNotFinalStatus(appCtx, g.selectLimit)
				if err != nil {
					logger.Log.Error("collect orders for update", zap.Error(err))
					continue
				}
				for _, order := range orders {
					select {
					case <-appCtx.Done(): //app exit signal
						break out
					case retryAfter = <-retryAfterChan: //retryAfter signal from poll workers
						//delay grabbing process for retryAfter minutes and return
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
