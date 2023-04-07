package locker

import (
	"errors"
	"time"
)

var errSmallTTL = errors.New("TTL must be bigger than refresh period")

type config struct {
	ttl           time.Duration
	refreshPeriod time.Duration
}

func (m *config) Clone() *config {
	return &config{
		ttl:           m.ttl,
		refreshPeriod: m.refreshPeriod,
	}
}

func (m *config) Validate() error {
	if m.ttl <= m.refreshPeriod {
		return errSmallTTL
	}

	return nil
}

func makeConfigWithOptions(cfg *config, options []Option) (*config, error) {
	cfg = cfg.Clone()
	for _, o := range options {
		o(cfg)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
