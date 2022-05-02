package vdisplay

import (
	"log"
	"os"
	"path/filepath"

	"github.com/inahga/vdisplay/internal/drm"
)

// VKMS uses the `vkms` kernel module.
//
// See https://github.com/torvalds/linux/tree/master/drivers/gpu/drm/vkms.
type VKMS struct {
	card *os.File
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
			fd, err := os.Open(p)
			if err != nil {
				log.Printf("vkms: open: %s", err)
				fd.Close()
				continue
			}

			ver, err := drm.CardVersion(fd)
			if err != nil {
				log.Printf("vkms: drm: %s", err)
				fd.Close()
				continue
			}

			if ver.Name != vkmsIdentifier {
				log.Printf("vkms: ignoring card %s because it is not vkms", p)
				fd.Close()
				continue
			}

			log.Printf("vkms: using card %s", p)
			found = true
			availableVDisplays = append(availableVDisplays, &VKMS{card: fd})
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
