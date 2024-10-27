package routine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/esafronov/yp-sysloyalty/internal/app/config"
	"github.com/esafronov/yp-sysloyalty/internal/domain"
	"github.com/esafronov/yp-sysloyalty/internal/logger"
	"go.uber.org/zap"
)

type ErrPollRetry struct {
	RetryAfter int
}

func (e *ErrPollRetry) Error() string {
	return fmt.Sprintf("retry after %d minutes", e.RetryAfter)
}

func NewErrPollRetry(retryAfter int) *ErrPollRetry {
	return &ErrPollRetry{
		RetryAfter: retryAfter,
	}
}

var ErrNoContent = errors.New("no content") //http.StatusNoContent

type Poller struct {
	workerCount    int                      //worker count
	RetryAfterChan chan int                 //channel for sending retryAfter signal
	resultChan     chan *domain.OrderUpdate //results channel
	delayed        bool                     //workers state
	endPoint       string                   //accrual system address
	workerTimeout  int                      //worker request timeout in milliseconds
}

// poller factory
func NewPoller(params *config.AppParams) *Poller {
	workerTimeout := 200                                                     //worker request timeout
	workerRate := float64(1000 / workerTimeout)                              //worker perfomance
	workerCount := int(math.Ceil(float64(*params.ProcessRate) / workerRate)) //calculating workers count needed
	return &Poller{
		workerCount:    workerCount,
		RetryAfterChan: make(chan int),
		resultChan:     make(chan *domain.OrderUpdate, 100),
		endPoint:       *params.AccrualSystemAddress,
		workerTimeout:  workerTimeout,
	}
}

// runs Poller routine and returns OrderUpdate channel
func (p *Poller) Run(ctx context.Context, orderChan <-chan *domain.Order, wg *sync.WaitGroup) chan *domain.OrderUpdate {
	wg.Add(1)
	var workerWg sync.WaitGroup //waitgroup for all poll workers
	for i := 1; i <= p.workerCount; i++ {
		go p.Worker(ctx, orderChan, &workerWg, i)
	}
	go func() {
		//wait until all poll workers are finished
		workerWg.Wait()
		close(p.resultChan)
		close(p.RetryAfterChan)
		logger.Log.Info("exit from poller...")
		wg.Done()
	}()
	return p.resultChan
}

// start poll worker
func (p *Poller) Worker(ctx context.Context, orderChan <-chan *domain.Order, wg *sync.WaitGroup, i int) {
	wg.Add(1)
	defer wg.Done()
	logger.Log.Info("poll worker started...", zap.Int("№", i))
	for order := range orderChan {
		select {
		case <-ctx.Done():
			return
		default:
			if p.delayed {
				continue
			}
			orderUpdate, err := p.requestUpdate(ctx, order.Num)
			if err != nil {
				if retErr, ok := err.(*ErrPollRetry); ok {
					p.RetryAfterChan <- retErr.RetryAfter
					p.delayed = true
					time.AfterFunc(time.Duration(retErr.RetryAfter), func() {
						p.delayed = false
					})
					logger.Log.Info("accrual system ask to retry", zap.String("order", order.Num), zap.Error(retErr))
					continue
				}
				if errors.Is(err, ErrNoContent) {
					logger.Log.Info("order is not found in accrual system", zap.String("order", order.Num), zap.Error(err))
					continue
				}
				logger.Log.Error("can't request from accrual system", zap.String("order", order.Num), zap.Error(err))
				continue
			}
			p.resultChan <- &orderUpdate
		}
	}
}

// request status for order from Accrual system
func (p *Poller) requestUpdate(ctx context.Context, orderNum string) (update domain.OrderUpdate, err error) {
	url := p.endPoint + "/api/orders/" + orderNum

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(p.workerTimeout)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, url, nil)
	if err != nil {
		err = fmt.Errorf("new request: %w", err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("do request: %w", err)
		return
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		err = json.NewDecoder(res.Body).Decode(&update)
	case http.StatusNoContent:
		err = ErrNoContent
	case http.StatusTooManyRequests:
		retryAfter, err := strconv.ParseInt(res.Header.Get("Retry-After"), 10, 64)
		if err == nil {
			err = NewErrPollRetry(int(retryAfter))
		}
	default:
		err = fmt.Errorf("response status: %d", res.StatusCode)
	}
	return
}
