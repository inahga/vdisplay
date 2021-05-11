// Package backend facilitates creation of virtual monitors and frame capture.
package backend

// Backend is a platform agnostic interface for provisioning virtual displays
// and capturing frames.
type Backend interface {
	// IsSupported returns whether the backend is supported on the current platform.
	IsSupported() bool
	// CreateVirtualDisplay sets up a new virtual display with the given parameters.
	CreateVirtualDisplay() error
}

// Best returns the backend that is best supported on the current platform.
func Best() (*Backend, error) {
	return nil, nil
}
