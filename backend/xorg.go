package backend

// Xorg uses Xrandr, Xfixes, Xshm, and Xlib to provision virtual screens and
// capture frames.
type Xorg struct {
}

func (x *Xorg) IsSupported() bool {
	// check if the current platform is linux
	// check if running in x mode
	// check for the presence of the xrandr binary (are we going to use that?)
	// check if driver supports the VIRTUAL displays

	return false
}
