//go:build go1.20
// +build go1.20

package locker

import "context"

var withCancelCause = context.WithCancelCause
