package ysword

import (
	"os"
	"syscall"

	"github.com/zlog-fun/ysword/graceful"
	"github.com/zlog-fun/ysword/server"
)

// AppInfoData is application context value.
type AppInfoData interface {
	Name() string
	Version() string
}

// App is an application components lifecycle manager.
type App struct {
	opts options
}

func New(opts ...Option) *App {
	o := options{
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		pid:  "pid",
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &App{
		opts: o,
	}
}

// Name returns service name.
func (a *App) Name() string { return a.opts.name }

// Version returns app version.
func (a *App) Version() string { return a.opts.version }

// Run start applicayion life cycle
func (a *App) Run() error {
	// startup http server
	srv := server.New(
		server.Host(a.opts.httpConfig.Host),
		server.Port(a.opts.httpConfig.Port),
		server.Debug(a.opts.httpConfig.Debug),
		server.Router(a.opts.handle),
	)
	http := srv.Run(a.Name())

	// auto switch restart component
	g := graceful.New(
		graceful.Host(a.opts.httpConfig.Host),
		graceful.Port(a.opts.httpConfig.Port),
		graceful.Method(a.opts.graceful.Method),
		graceful.Pid(a.opts.graceful.Pid),
		graceful.Handler(http),
	)
	g.Bootstrap()

	return nil
}
