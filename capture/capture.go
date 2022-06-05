// Package capture provides interfaces for capturing virtual display output.
package capture

import "image"

type Buffer struct{}

type Capture interface {
	Close() error
	Register(func(image.Image))
	SetMaxFramerate(uint32)
	Start() error
}
