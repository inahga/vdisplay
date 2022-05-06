// Package capture provides interfaces for capturing virtual display output.
package capture

type Buffer struct{}

type Capture interface {
	Close() error
	Register(chan<- Buffer)
	SetMaxFramerate(uint32)
	Start() error
}
