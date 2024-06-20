package logger

import (
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	azlog "github.com/jensneuse/abstractlogger"
)

type (
	Logger struct {
		Log zerolog.Logger
	}
)

func (g Logger) Enabled(level zapcore.Level) bool {
	switch level {
	case zapcore.DebugLevel:
		return g.Log.Debug().Enabled()
	case zapcore.DPanicLevel:
		return g.Log.Panic().Enabled()
	case zapcore.ErrorLevel:
		return g.Log.Error().Enabled()
	case zapcore.WarnLevel:
		return g.Log.Warn().Enabled()
	case zapcore.FatalLevel:
		return g.Log.Fatal().Enabled()
	case zapcore.InfoLevel:
		return g.Log.Info().Enabled()
	case zapcore.PanicLevel:
		return g.Log.Panic().Enabled()
	case zapcore.InvalidLevel:
		return false
	}
	return false
}

func (g Logger) With(fields []zapcore.Field) zapcore.Core {
	fieldMap := map[string]any{}
	for _, field := range fields {
		fieldMap[field.Key] = field.Interface
	}
	g.Log = g.Log.With().Fields(fieldMap).Logger()
	return g
}

func (g Logger) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if g.Enabled(entry.Level) {
		return ce.AddCore(entry, g)
	}
	return ce
}

func (g Logger) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	fieldMap := map[string]any{}
	var err error
	for _, field := range fields {
		switch v := field.Interface.(type) {
		case error:
			err = v
		default:
			fieldMap[field.Key] = field.Interface
		}
	}
	switch entry.Level {
	case zapcore.DebugLevel:
		g.Log.Debug().Fields(fieldMap).Msgf(entry.Message)
	case zapcore.DPanicLevel, zapcore.PanicLevel:
		if err != nil {
			g.Log.Panic().Err(err).Fields(fieldMap).Msgf(entry.Message)
		} else {
			g.Log.Panic().Fields(fieldMap).Msgf(entry.Message)
		}
	case zapcore.ErrorLevel:
		if err != nil {
			g.Log.Error().Err(err).Fields(fieldMap).Msgf(entry.Message)
		} else {
			g.Log.Error().Fields(fieldMap).Msgf(entry.Message)
		}
	case zapcore.WarnLevel:
		g.Log.Warn().Fields(fieldMap).Msgf(entry.Message)
	case zapcore.FatalLevel:
		if err != nil {
			g.Log.Fatal().Err(err).Fields(fieldMap).Msgf(entry.Message)
		} else {
			g.Log.Fatal().Fields(fieldMap).Msgf(entry.Message)
		}
	case zapcore.InfoLevel:
		g.Log.Info().Fields(fieldMap).Msgf(entry.Message)
	case zapcore.InvalidLevel:
		// ignored
	}

	return nil
}

func (g Logger) Sync() error {
	return nil
}

func ZapLogger(log zerolog.Logger) azlog.Logger {
	logger, err := zap.NewDevelopmentConfig().Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &Logger{Log: log}
	}))
	if err != nil {
		panic(err)
	}

	return azlog.NewZapLogger(logger, azlog.DebugLevel)
}
