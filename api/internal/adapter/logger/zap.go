package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	src "microservice"
	"microservice/config"
	"microservice/internal/adapter/registry"
	"os"
	"time"
)

type level int

const (
	_ level = iota
	DEBUG
	INFO
	WARN
	ERR
	FATAL
)

type logger struct {
	service config.Service
	config  config.Logger
	zap     *zap.Logger
}

func New(registry registry.IRegistry) ILogger {
	client := new(logger)
	registry.Parse(&client.service)
	registry.Parse(&client.config)

	return client
}

func (l *logger) Init() {
	var opts []zap.Option

	core := zapcore.NewTee(
		consoleCore(),
		fileCore(ERR, "err"),
		fileCore(INFO, "info"),
	)

	opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1))

	l.zap = zap.New(core, opts...).With(zap.String("usecase", l.service.Env))
}

func (l *logger) Stop() {
	_ = l.zap.Sync() // ignore Sync error: because the stdout isn't flushable
	log.Printf("[logger] zap stopped successfully")
}

func (l *logger) C() *zap.Logger                          { return l.zap }
func (l *logger) Debug(scope string, fields ...zap.Field) { l.zap.Debug(scope, fields...) }
func (l *logger) Info(scope string, fields ...zap.Field)  { l.zap.Info(scope, fields...) }
func (l *logger) Warn(scope string, fields ...zap.Field)  { l.zap.Warn(scope, fields...) }
func (l *logger) Error(scope string, fields ...zap.Field) { l.zap.Error(scope, fields...) }

// HELPERS

func levelEnabler(lvl level) zap.LevelEnablerFunc {
	levels := map[level]zap.LevelEnablerFunc{
		DEBUG: zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl >= zapcore.DebugLevel }), // Logs everything to stdout
		INFO:  zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl == zapcore.InfoLevel }),
		WARN:  zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl == zapcore.WarnLevel }),
		ERR:   zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl == zapcore.ErrorLevel }),
	}

	return levels[lvl]
}

func consoleCore() zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		levelEnabler(DEBUG),
	)
}

func fileCore(lvl level, prefix string) zapcore.Core {
	fileName := fmt.Sprintf("%s/logs/%s-%s.log", src.Root(), prefix, time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: there could be multiple Cores per level

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(file),
		levelEnabler(lvl),
	)
}
