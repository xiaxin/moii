package log

import (
	"fmt"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	skip = 1
	log  = New(
		nil,
		zap.AddCallerSkip(skip),
		zap.AddStacktrace(zap.ErrorLevel),
	)
)

const (
	_oddNumberErrMsg    = "Ignored key without a value."
	_nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
)

const (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
)

type Logger struct {
	base *zap.Logger
}

func New(conf *zap.Config, opts ...zap.Option) *Logger {
	config := &zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		// use console or json
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			// time
			TimeKey: "T",
			// 默认 zapcore.ISO8601TimeEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			},

			// level
			LevelKey:    "L",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			CallerKey:    "C",
			EncodeCaller: zapcore.ShortCallerEncoder,

			NameKey:        "N",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
		//  TODO
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// TODO 适配
	if nil != conf {
		if nil != conf.OutputPaths {
			config.OutputPaths = conf.OutputPaths
		}

		if nil != conf.ErrorOutputPaths {
			config.ErrorOutputPaths = conf.ErrorOutputPaths
		}
	}

	//  Logger
	logger, _ := config.Build(zap.AddCallerSkip(2))

	return &Logger{
		base: logger.WithOptions(opts...),
	}
}


func (l *Logger) Named(name string) *Logger {
	return &Logger{
		base: l.base.Named(name),
	}
}

func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		base: l.base.With(l.sweetenFields(args)...),
	}
}

func (l *Logger) Clone(opts ...zap.Option) *Logger {
	return &Logger{
		base: l.base.WithOptions(opts...),
	}
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(zap.DebugLevel, "", args, nil)
}

func (l *Logger) Info(args ...interface{}) {
	l.log(zap.InfoLevel, "", args, nil)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log(zap.WarnLevel, "", args, nil)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *Logger) Error(args ...interface{}) {
	l.log(zap.ErrorLevel, "", args, nil)
}

func (l *Logger) DPanic(args ...interface{}) {
	l.log(zap.PanicLevel, "", args, nil)
}

func (l *Logger) Panic(args ...interface{}) {
	l.log(zap.PanicLevel, "", args, nil)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log(zap.FatalLevel, "", args, nil)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.log(zap.DebugLevel, template, args, nil)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.log(zap.InfoLevel, template, args, nil)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.log(zap.WarnLevel, template, args, nil)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.log(zap.ErrorLevel, template, args, nil)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.log(zap.PanicLevel, template, args, nil)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.log(zap.PanicLevel, template, args, nil)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.log(zap.FatalLevel, template, args, nil)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.log(zap.DebugLevel, msg, nil, keysAndValues)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.log(zap.InfoLevel, msg, nil, keysAndValues)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.log(zap.WarnLevel, msg, nil, keysAndValues)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.log(zap.ErrorLevel, msg, nil, keysAndValues)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.log(zap.DPanicLevel, msg, nil, keysAndValues)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.log(zap.PanicLevel, msg, nil, keysAndValues)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.log(zap.FatalLevel, msg, nil, keysAndValues)
}

func (l *Logger) Sync() error {
	return l.base.Sync()
}

//  统一日志记录类
func (l *Logger) log(lvl zapcore.Level, template string, fmtArgs []interface{}, context []interface{}) {
	if lvl < zap.DPanicLevel && !l.base.Core().Enabled(lvl) {
		return
	}

	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	if ce := l.base.Check(lvl, msg); ce != nil {

		context := l.sweetenFields(context)
		ce.Write(context...)
	}
}

func (l *Logger) sweetenFields(args []interface{}) []zapcore.Field {
	if len(args) == 0 {
		return nil
	}

	fields := make([]zapcore.Field, 0, len(args))
	var invalid invalidPairs

	for i := 0; i < len(args); {
		if f, ok := args[i].(zap.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		if i == len(args)-1 {
			l.base.DPanic(_oddNumberErrMsg, zap.Any("ignored", args[i]))
			break
		}

		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			if cap(invalid) == 0 {
				invalid = make(invalidPairs, 0, len(args)/2)
			}
			invalid = append(invalid, invalidPair{i, key, val})
		} else {
			fields = append(fields, zap.Any(keyStr, val))
		}
		i += 2
	}

	if len(invalid) > 0 {
		l.base.DPanic(_nonStringKeyErrMsg, zap.Array("invalid", invalid))
	}
	return fields
}

type invalidPair struct {
	position   int
	key, value interface{}
}

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("position", int64(p.position))
	zap.Any("key", p.key).AddTo(enc)
	zap.Any("value", p.value).AddTo(enc)
	return nil
}

type invalidPairs []invalidPair

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	for i := range ps {
		err = multierr.Append(err, enc.AppendObject(ps[i]))
	}
	return err
}

// default
func Named(name string) *Logger {
	log = log.Named(name)
	return log
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(template string, args ...interface{}) {
	log.Infof(template, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	log.Warnf(template, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args)
}

func Fatalf(template string, args ...interface{}) {
	log.Fatalf(template, args...)
}

func Clone() *Logger {
	return log.Clone()
}
