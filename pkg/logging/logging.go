package logging

import (
	"context"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	encoder zapcore.Encoder
	ws      zapcore.WriteSyncer
	level   zapcore.LevelEnabler
}

func NewLogger() *Logger {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	})
	writeStream := zapcore.Lock(zapcore.AddSync(os.Stdout))
	errorStream := zap.ErrorOutput(zapcore.Lock(zapcore.AddSync(os.Stderr)))
	callerOption := zap.WithCaller(true)
	level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.ErrorLevel
	})
	return &Logger{
		Logger: zap.New(zapcore.NewCore(encoder, writeStream, level)).
			WithOptions(
				callerOption,
				errorStream,
			),
		encoder: encoder,
		ws:      writeStream,
		level:   level,
	}
}

func (l *Logger) WithField(field Field) *Logger {
	return &Logger{
		Logger:  l.With(zap.Any(field.Label, field.Value)),
		encoder: l.encoder,
		ws:      l.ws,
		level:   l.level,
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	if val, ok := ctx.Value(Logger{}).([]interface{}); ok && val != nil {
		return &Logger{
			Logger:  l.Sugar().With(val...).Desugar(),
			encoder: l.encoder,
			ws:      l.ws,
			level:   l.level,
		}
	}

	return l
}

func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger:  l.With(zap.Error(err)),
		encoder: l.encoder,
		ws:      l.ws,
		level:   l.level,
	}
}

func (l *Logger) SetOutput(out io.Writer) {
	l.ws = zapcore.Lock(zapcore.AddSync(out))
	l.Logger = zap.New(zapcore.NewCore(l.encoder, l.ws, l.level)).WithOptions(zap.WithCaller(true))
}

func (l *Logger) SetErrorOutput(out io.Writer) {
	l.Logger = l.Logger.WithOptions(zap.ErrorOutput(zapcore.Lock(zapcore.AddSync(out))))
}

var DefaultLogger = NewLogger()

func InitDefaultLogger(ctx context.Context) {
	DefaultLogger = DefaultLogger.WithContext(ctx)
}

func Info(msg string) {
	DefaultLogger.Info(msg)
}

func Error(msg string) {
	DefaultLogger.Error(msg)
}

func Fatal(msg string) {
	DefaultLogger.Fatal(msg)
}

func Debug(msg string) {
	DefaultLogger.Debug(msg)
}

func WithField(field Field) *Logger {
	return DefaultLogger.WithField(field)
}

func WithError(err error) *Logger {
	return DefaultLogger.WithError(err)
}

func WithContext(ctx context.Context) *Logger {
	return DefaultLogger.WithContext(ctx)
}

func SetOutput(out io.Writer) {
	DefaultLogger.SetOutput(out)
}

func SetErrorOutput(out io.Writer) {
	DefaultLogger.SetErrorOutput(out)
}
