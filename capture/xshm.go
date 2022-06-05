//go:build cgo && (linux || freebsd || openbsd || dragonfly)

package capture

// Xshm uses X shared memory extensions for capturing output. It only makes sense
// if X is the display server.
type Xshm struct{}

// See https://www.ssec.wisc.edu/~billh/bp/xshm.c for a minimal implementation
// of this strategy.
