# vkms
A DRM kernel module for creating an in-memory framebuffer.
This is a fork of the in-tree kernel module `vkvm` (`drivers/gpu/drm/vkvm`).
See https://dri.freedesktop.org/docs/drm/gpu/vkms.html.

This module is highly unstable when run alongside a regular GPU.
Since it's out of tree, it'll also taint your kernel.
At the moment, use should only occur in a VM.

### Local Dev Setup
To build and install the kernel module (assuming Fedora):
1. Ensure your environment is set up for building kernel modules. `sudo dnf builddep kernel` should do the trick.
1. Use `make` to build `vkms.ko`.
1. Execute the following commands:
```bash
sudo cp vkms.ko /lib/modules/$(uname -r)
sudo depmod -a
echo 'vkms' | sudo tee /etc/modules-load.d/vkms.conf
```
1. Reboot the system.

You can hot-add it with `insmod`, but your desktop environment and/or kernel will almost certainly crash.
