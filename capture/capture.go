// Package capture provides interfaces for capturing virtual display output.
package capture

type Buffer struct{}

type Capture interface {
	Close() error
	Register(func(buf []byte))
	SetMaxFramerate(uint32)
	Start() error
}
