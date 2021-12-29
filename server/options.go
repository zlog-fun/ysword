package server

import "github.com/gin-gonic/gin"

// Option is an application option.
type Option func(o *options)

type options struct {
	host  string
	port  int64
	debug string
	// TODO
	handle func(*gin.Engine)
}

func Host(host string) Option {
	return func(o *options) { o.host = host }
}

func Port(port int64) Option {
	return func(o *options) { o.port = port }
}

func Debug(debug string) Option {
	return func(o *options) { o.debug = debug }
}

func Router(callback func(engine *gin.Engine)) Option {
	return func(o *options) { o.handle = callback }
}
