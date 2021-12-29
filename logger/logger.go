package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var lg *zap.Logger

type Logger struct {
	opts options
}

func (log *Logger) InitLogger() (err error) {
	writeSyncer := getLogWriter(log.opts.filename, log.opts.maxsize, log.opts.maxbackups, log.opts.maxage)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(log.opts.level))
	if err != nil {
		return
	}
	if log.opts.level == "debug" {
		// output os.Stdout
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)
	development := zap.Development()
	lg = zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel), // Error 级别增加调用堆栈
		development,                       // 增加行号
	)
	// 替换zap包中全局的logger实例
	zap.ReplaceGlobals(lg)
	// 初始化后可以直接 zap.L() 获取 *Logger
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	// call line and code path
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func New(opts ...Option) *Logger {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &Logger{
		opts: o,
	}
}
