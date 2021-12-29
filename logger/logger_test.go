package logger

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestZap(t *testing.T) {
	logger := zap.NewExample()
	defer logger.Sync()
	logger.Info("failed to fetch URL",
		zap.String("url", "http://example.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	sugar.Infow("failed to fetch URL",
		"url", "http://example.com",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("failed to fetch URL: %s", "http://example.com")
}

func TestLogger(t *testing.T) {
	var config = New(
		Level("denug"),
		Filename("1.log"),
		MaxAge(1),
		Maxbackups(1),
		MaxSize(10),
	)
	config.InitLogger()
	zap.L().Error("hello")
}
