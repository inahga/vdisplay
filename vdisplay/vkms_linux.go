package vdisplay

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/inahga/vdisplay/capture"
	"github.com/inahga/vdisplay/internal/drm"
)

// VKMS uses the `vkms` kernel module.
//
// See https://github.com/torvalds/linux/tree/master/drivers/gpu/drm/vkms.
type VKMS struct {
	card    *drm.Card
	capture capture.Capture
}

const (
	vkmsDRIDir     = "/dev/dri"
	vkmsIdentifier = "vkms"
)

func newVKMS(c *drm.Card) (*VKMS, error) {
	ret := &VKMS{card: c}
	ver, err := c.Version()
	if err != nil {
		return nil, err
	}
	if ver.Name != vkmsIdentifier {
		return nil, fmt.Errorf("card is not vkms")
	}
	return ret, nil
}

func init() {
	files, err := os.ReadDir(vkmsDRIDir)
	if err != nil {
		log.Printf("[vkms] readdir: %s", files)
	}

	var found bool
	for _, f := range files {
		if f.Type()&os.ModeDevice != 0 {
			p := filepath.Join(vkmsDRIDir, f.Name())
			c, err := drm.Open(p)
			if err != nil {
				log.Printf("[vkms] open: %s", err)
				continue
			}
			vkms, err := newVKMS(c)
			if err != nil {
				log.Printf("[vkms] %s: %s", p, err)
				c.Close()
				continue
			}
			log.Printf("[vkms] using card %s", p)
			found = true
			availableVDisplays = append(availableVDisplays, vkms)
			return
		}
	}
	if !found {
		log.Printf("[vkms] no cards found, is the kernel module enabled?")
	}
}

func (v *VKMS) priority() int {
	return 100
}
