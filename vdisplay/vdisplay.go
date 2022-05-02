// Package vdisplay manages virtual displays according to what the platform
// supports.
package vdisplay

// VDisplay is a virtual display handler.
type VDisplay interface {
	priority() int
}

var availableVDisplays []VDisplay
