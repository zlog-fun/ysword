package logger

// Option is an application option.
type Option func(o *options)

type options struct {
	level      string `yaml:"level"`
	filename   string `yaml:"filename"`
	maxsize    int    `yaml:"maxsize"`
	maxage     int    `yaml:"max_age"`
	maxbackups int    `yaml:"max_backups"`
}

func Level(level string) Option {
	return func(o *options) { o.level = level }
}

func Filename(filename string) Option {
	return func(o *options) { o.filename = filename }
}

func MaxSize(maxsize int) Option {
	return func(o *options) { o.maxsize = maxsize }
}

func MaxAge(maxage int) Option {
	return func(o *options) { o.maxage = maxage }
}

func Maxbackups(maxbackups int) Option {
	return func(o *options) { o.maxbackups = maxbackups }
}
