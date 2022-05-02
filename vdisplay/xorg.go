//go:build cgo && (linux || freebsd || openbsd || dragonfly)

package vdisplay

// Xorg sets up a VIRTUAL display with xrandr, and captures screenshots with
// Xlib. Not all drivers support a VIRTUAL display.
type Xorg struct{}

func (x *Xorg) priority() int {
	return 25
}
