package server

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HttpConfig struct {
	Host  string `yaml:"host"`
	Port  int64  `yaml:"port"`
	Debug string `yaml:"debug"`
}

// Http [http server struct]
type Http struct {
	opts options
}

type Register func(engine *gin.Engine)

// Run register is gin handler
func (s Http) Run(appname string) *gin.Engine {
	// setting run mode
	if s.opts.debug != "" {
		gin.SetMode(s.opts.debug)
	}
	engine := gin.New()

	// load  middleware
	engine.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	// Logs all panic to response log
	//   - stack means whether output the stack info.
	engine.Use(ginzap.RecoveryWithZap(zap.L(), true))

	// register gin engine
	s.opts.handle(engine)

	return engine
}

func New(opts ...Option) *Http {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &Http{
		opts: o,
	}
}
