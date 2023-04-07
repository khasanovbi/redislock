//go:build go1.20
// +build go1.20

package redislock

import "context"

var withCancelCause = context.WithCancelCause
