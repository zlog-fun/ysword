package config

type Option func(*options)

type options struct {
	name string
	path string
	env  string
}

func Name(name string) Option {
	return func(o *options) { o.name = name }
}

func Path(path string) Option {
	return func(o *options) { o.path = path }
}

// load different config
func Env(env string) Option {
	return func(o *options) { o.env = env }
}
