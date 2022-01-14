package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lv = zap.NewAtomicLevel()

func init() {
	defaultConfig := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:   "T",
		LevelKey:  "L",
		NameKey:   "N",
		CallerKey: "C",

		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   nil,
	}
	consoleEncoder := zapcore.NewConsoleEncoder(defaultConfig)

	stdoutSyncer := zapcore.Lock(os.Stdout)
	// tee core
	cfgcore := zapcore.NewTee(
		zapcore.NewCore(
			consoleEncoder,
			stdoutSyncer,
			lv,
		),
	)
	logger := zap.New(cfgcore, zap.AddCaller())
	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)
	logger.Sync()
}

func ChangeLogLevel(s string) {
	lv.UnmarshalText([]byte(s))
}
