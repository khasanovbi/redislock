package redislock

import (
	"time"
)

var defaultConfig = &config{
	ttl:           10 * time.Minute, //nolint:gomnd
	refreshPeriod: time.Minute,
}

type Option func(m *config)

func WithTTL(ttl time.Duration) Option {
	return func(m *config) {
		m.ttl = ttl
	}
}

func WithRefreshPeriod(refreshPeriod time.Duration) Option {
	return func(m *config) {
		m.refreshPeriod = refreshPeriod
	}
}
