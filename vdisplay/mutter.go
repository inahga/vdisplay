//go:build linux || freebsd || openbsd || dragonfly

package vdisplay

// Mutter uses the GNOME Window Manager.
//
// Creation of a virtual display is supported in GNOME 40.
// See https://gitlab.gnome.org/GNOME/mutter/-/merge_requests/1698.
type Mutter struct{}

func (m *Mutter) priority() int {
	return 50
}
