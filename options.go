package app

import (
	"context"
	"os"

	"ysword/graceful"
	"ysword/server"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ServerEngine interface {
	Run(interface{})
}

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	ctx        context.Context
	sigs       []os.Signal
	name       string
	version    string
	pid        string
	logger     *zap.Logger
	handle     func(engine *gin.Engine)
	httpConfig server.HttpConfig
	graceful   graceful.GracefulConfig
}

// Name with service version.
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version with service version.
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Pid file path
func Pid(pid string) Option {
	return func(o *options) { o.pid = pid }
}

// Context with service context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Logger with service logger.
func Logger(logger *zap.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func Handle(f func(engine *gin.Engine)) Option {
	return func(o *options) {
		o.handle = f
	}
}

// server config
func HttpConfig(c server.HttpConfig) Option {
	return func(o *options) {
		o.httpConfig = c
	}
}

func Graceful(g graceful.GracefulConfig) Option {
	return func(o *options) {
		o.graceful = g
	}
}
