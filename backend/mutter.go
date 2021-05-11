package backend

// Mutter uses Mutter, the GNOME Window Manager, to provision virtual displays and
// capture frames.
//
// Creation of a virtual display is supported in GNOME 40.
// See https://gitlab.gnome.org/GNOME/mutter/-/merge_requests/1698.
type Mutter struct {
}

func (m *Mutter) IsSupported() bool {
	// is wayland only supported?

	// check if the current platform is linux
	// check if gnome-shell is running
	// check the version of gnome-shell
	// check if running in wayland or xorg (can use loginctl) ???

	// this backend is broken.
	// follow https://gitlab.gnome.org/GNOME/mutter/-/issues/1728
	return false
}
