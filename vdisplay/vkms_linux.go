package vdisplay

import (
	"log"
	"os"
	"path/filepath"

	"github.com/inahga/vdisplay/drm"
)

// VKMS uses the `vkms` kernel module.
//
// See https://github.com/torvalds/linux/tree/master/drivers/gpu/drm/vkms.
type VKMS struct {
	card *drm.Card
}

const (
	vkmsDRIDir     = "/dev/dri"
	vkmsIdentifier = "vkms"
)

func init() {
	files, err := os.ReadDir(vkmsDRIDir)
	if err != nil {
		log.Printf("vkms: readdir: %s", files)
	}

	var found bool
	for _, f := range files {
		if f.Type()&os.ModeDevice != 0 {
			p := filepath.Join(vkmsDRIDir, f.Name())
			c, err := drm.Open(p)
			if err != nil {
				log.Printf("vkms: open: %s", err)
				continue
			}

			ver, err := c.Version()
			if err != nil {
				log.Printf("vkms: drm: %s", err)
				c.Close()
				continue
			}

			if ver.Name != vkmsIdentifier {
				log.Printf("vkms: ignoring card %s because it is not vkms", p)
				c.Close()
				continue
			}

			if err := c.SetClientCap(drm.ClientCapAtomic, 1); err != nil {
				log.Printf("vkms: %s: setcap atomic: %s", p, err)
				c.Close()
				continue
			}

			if err := c.SetClientCap(drm.ClientCapWritebackConnectors, 1); err != nil {
				log.Printf("vkms: %s: setcap writeback: %s", p, err)
				c.Close()
				continue
			}

			log.Printf("vkms: using card %s: %+v", p, ver)
			found = true
			availableVDisplays = append(availableVDisplays, &VKMS{card: c})
			return
		}
	}
	if !found {
		log.Printf("vkms: no cards found, is the kernel module enabled?")
	}
}

func (v *VKMS) priority() int {
	return 100
}
