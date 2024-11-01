package main

import (
	"errors"
	"net/http"

	"github.com/esafronov/yp-sysloyalty/internal/app"
	"github.com/esafronov/yp-sysloyalty/internal/logger"
	"go.uber.org/zap"
)

func main() {
	err := logger.Initialize("debug")
	if err != nil {
		panic("can't initialize logger")
	}
	if err := app.Run(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("server shutdown", zap.Error(err))
		} else {
			logger.Log.Error("app error", zap.Error(err))
		}
	}
}
