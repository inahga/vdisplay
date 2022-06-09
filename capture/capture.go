// Package capture provides interfaces for capturing virtual display output.
package capture

import "image"

type Capture interface {
	Close() error
	Start(framerate uint32, rect image.Rectangle, cb func(image.Image)) error
}
