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

var ErrNoContent = errors.New("no content")

type Poller struct {
	workerCount    int
	RetryAfterChan chan int
	resultChan     chan *domain.OrderUpdate
	delayed        bool
	endPoint       string
	workerTimeout  int //milliseconds
}

func NewPoller(params *config.AppParams) *Poller {
	workerTimeout := 200
	workerRate := float64(1000 / workerTimeout)
	workerCount := int(math.Ceil(float64(*params.ProcessRate) / workerRate))
	return &Poller{
		workerCount:    workerCount,
		RetryAfterChan: make(chan int),
		resultChan:     make(chan *domain.OrderUpdate, 100),
		endPoint:       *params.AccrualSystemAddress,
		workerTimeout:  workerTimeout,
	}
}

func (p *Poller) Run(ctx context.Context, orderChan <-chan *domain.Order, wg *sync.WaitGroup) chan *domain.OrderUpdate {
	wg.Add(1)
	var workerWg sync.WaitGroup
	for i := 1; i <= p.workerCount; i++ {
		go p.Worker(ctx, orderChan, &workerWg, i)
	}
	go func() {
		workerWg.Wait()
		close(p.resultChan)
		close(p.RetryAfterChan)
		logger.Log.Info("exit from poller...")
		wg.Done()
	}()
	return p.resultChan
}

func (p *Poller) Worker(ctx context.Context, orderChan <-chan *domain.Order, wg *sync.WaitGroup, i int) {
	wg.Add(1)
	defer wg.Done()
	logger.Log.Info("poll worker started...", zap.Int("num", i))
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
