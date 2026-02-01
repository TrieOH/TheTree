package apierr

import (
	"GoAuth/internal/adapters/observability/logs"

	"github.com/MintzyG/fail"
	"go.uber.org/zap"
)

func OnFromFailHook(e error) {
	logs.L().Info("from hook fail", zap.Any("generic error", e))
}

func OnFromSuccessHook(e error, fe *fail.Error) {
	logs.L().Info("from hook success", zap.Any("generic error", e), zap.Any("transformed error", fe), zap.Any("transformed error dump", fe.Dump()))
}
