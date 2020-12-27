// Package drm provides functions and structs for interacting with the Linux
// kernel's Direct Rendering Manager. It is primarily based off of the drm/drm_*.h
// headers.
package drm

import (
	"bytes"
	"fmt"
	"os"
	"unsafe"
)

type Card struct {
	fd *os.File
}

func New(fd *os.File) *Card {
	return &Card{fd: fd}
}

func Open(path string) (*Card, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &Card{
		fd: f,
	}, nil
}

func (c *Card) Close() error {
	return c.fd.Close()
}

func cToGoString(b []byte) string {
	return string(bytes.Trim(b, "\u0000"))
}

func (c *Card) Version() (*Version, error) {
	var ver cVersion
	if err := ioctl(c.fd, ioctlVersion, uintptr(unsafe.Pointer(&ver))); err != nil {
		return nil, fmt.Errorf("ioctl: %w", err)
	}

	var name, date, desc []byte
	if ver.namelen > 0 {
		name = make([]byte, ver.namelen+1)
		ver.name = uint64(uintptr(unsafe.Pointer(&name[0])))
	}
	if ver.datelen > 0 {
		date = make([]byte, ver.datelen+1)
		ver.date = uint64(uintptr(unsafe.Pointer(&date[0])))
	}
	if ver.desclen > 0 {
		desc = make([]byte, ver.desclen+1)
		ver.desc = uint64(uintptr(unsafe.Pointer(&desc[0])))
	}

	if err := ioctl(c.fd, ioctlVersion, uintptr(unsafe.Pointer(&ver))); err != nil {
		return nil, fmt.Errorf("ioctl: %w", err)
	}
	return &Version{
		Major:      ver.major,
		Minor:      ver.minor,
		PatchLevel: ver.patchlevel,
		Name:       cToGoString(name[:len(name)-1]),
		Date:       cToGoString(date[:len(date)-1]),
		Desc:       cToGoString(desc[:len(desc)-1]),
	}, nil
}

func (c *Card) SetClientCap(cap uint64, val uint64) error {
	setcap := cSetClientCap{
		capability: cap,
		value:      val,
	}
	if err := ioctl(c.fd, ioctlSetClientCap, uintptr(unsafe.Pointer(&setcap))); err != nil {
		return fmt.Errorf("ioctl: %w", err)
	}
	return nil
}

func (c *Card) SetMaster() error {
	return ioctl(c.fd, ioctlSetMaster, 0)
}

func (c *Card) DropMaster() error {
	return ioctl(c.fd, ioctlDropMaster, 0)
}
