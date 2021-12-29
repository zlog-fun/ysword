package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config is interface
type Config interface {
	Load() error
	Watch(v interface{}) error
	Scan(v interface{}) error
	Close() error
}

type config struct {
	opts options
}

func New(opts ...Option) Config {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &config{opts: o}
}

// Load init viper
func (c *config) Load() error {
	viper.SetConfigType("yaml")
	// according to env load different config file
	if c.opts.env != "" {
		viper.SetConfigName("config_" + c.opts.env)
	} else {
		viper.SetConfigName("config")
	}
	viper.AddConfigPath(c.opts.path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (c *config) Scan(v interface{}) error {
	err := viper.Unmarshal(v)
	if err != nil {
		return err
	}
	c.Watch(v)
	return nil
}

func (c *config) Watch(v interface{}) error {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		err := viper.Unmarshal(v)
		if err != nil {
			fmt.Printf("conf reload err: %s \n", err.Error())
		}
	})
	return nil
}

func (c *config) Close() error {
	return nil
}
