package graceful

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
)

type GracefulConfig struct {
	Method string `yaml:"method"`
	Pid    string `yaml:"pid"`
}

type Graceful struct {
	opts options
}

func (g *Graceful) Bootstrap() {
	sysType := runtime.GOOS
	if g.opts.method == "overseer" && sysType == "linux" {
		g.OverseerGraceful()
	} else {
		log.Println("Server starting. is DefaultGraceful")
		g.DefaultGraceful()
	}
}

// DefaultGraceful support win or linux，reload pid will change
func (g *Graceful) DefaultGraceful() {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", g.opts.host, g.opts.port),
		Handler: g.opts.handler,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 1 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

// OverseerGraceful is use overseer libary, pid don't change
func (g *Graceful) OverseerGraceful() {
	overseer.Run(overseer.Config{
		Program: g.prog,
		Address: fmt.Sprintf("%s:%d", g.opts.host, g.opts.port),
		Debug:   false,
		Fetcher: &fetcher.File{
			Path:     g.opts.pid,
			Interval: 3 * time.Second,
		},
	})
}

func (g *Graceful) prog(state overseer.State) {
	http.Serve(state.Listener, g.opts.handler)
}

func New(opts ...Option) *Graceful {
	o := options{
		pid: "pid",
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &Graceful{
		opts: o,
	}
}
