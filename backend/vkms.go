package backend

// VKMS uses the VKMS kernel module for provisioning virtual displays and capturing
// frames.
//
// It can be found at https://github.com/torvalds/linux/tree/master/drivers/gpu/drm/vkms.
//
// Note that some distributions' kernels aren't built with this kernel module.
// In that case, the module need to be installed manually, causing a kernel taint.
type VKMS struct {
}

func (v *VKMS) IsSupported() bool {
	// check if the current platform is linux
	// check for presence of the vkms kernel module

	return false
}
