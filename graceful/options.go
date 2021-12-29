package graceful

import "net/http"

type Option func(*options)

type options struct {
	method  string
	pid     string
	host    string
	port    int64
	handler http.Handler
}

func Method(method string) Option {
	return func(o *options) { o.method = method }
}

func Host(host string) Option {
	return func(o *options) { o.host = host }
}

func Pid(pid string) Option {
	return func(o *options) { o.host = pid }
}

func Port(port int64) Option {
	return func(o *options) { o.port = port }
}

func Handler(srv http.Handler) Option {
	return func(o *options) { o.handler = srv }
}
