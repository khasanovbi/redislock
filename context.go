//go:build !go1.20
// +build !go1.20

package locker

import "context"

func withCancelCause(ctx context.Context) (context.Context, cancelCauseFunc) {
	ctx, cancel := context.WithCancel(ctx)

	return ctx, func(cause error) {
		cancel()
	}
}
