package app

import (
	"context"

	"lib/errx"
	"lib/telemetry"
	"lib/validator"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type SimpleLogger struct{}

func (l *SimpleLogger) Intercept(_ context.Context, rs *fun.Response, statusCode int) {
	if statusCode == 500 {
		telemetry.Log().Info("InternalServerError Response", zap.Any("response", rs))
	}
}

func (l *SimpleLogger) InterceptSimple(rs *fun.Response, statusCode int) {
	if statusCode == 500 {
		telemetry.Log().Info("InternalServerError Response", zap.Any("response", rs))
	}
}

func SetupFUN() {
	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        app.cfg.AppName,
	})

	err := fun.AddInterceptor(&SimpleLogger{})
	if err != nil {
		errx.Exit(err, "failed to add interceptor")
	}

	v := validator.SetupValidator()
	bind.SetValidator(v)
	fun.SetPathParamFunc(chi.URLParam)
}
