package logger

import (
	"go.uber.org/zap"
)

//go:generate mockgen -source=./contract.go -destination=./mocks/logger_mock.go -package=logger_mock
type ILogger interface {
	Init()
	Stop()
	// C logger client
	C() *zap.Logger
	Debug(scope string, fields ...zap.Field)
	Info(scope string, fields ...zap.Field)
	Warn(scope string, fields ...zap.Field)
	Error(scope string, fields ...zap.Field)
}
