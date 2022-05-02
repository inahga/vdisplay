package drm

import (
	"bytes"
	"fmt"
	"os"
	"unsafe"
)

type (
	kernelSize = uint64
	cint       = int32
)

type cVersion struct {
	major      cint
	minor      cint
	patchlevel cint
	namelen    kernelSize
	name       uint64 // ptr to a []byte
	datelen    kernelSize
	date       uint64 // ptr to a []byte
	desclen    kernelSize
	desc       uint64 // ptr to a []byte
}

type Version struct {
	Major      int32
	Minor      int32
	PatchLevel int32
	Name       string
	Date       string
	Desc       string
}

var (
	ioctlVersion = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cVersion{})), ioctlBase, 0x00)
)

func CardVersion(c *os.File) (*Version, error) {
	var ver cVersion
	if err := ioctl(c, ioctlVersion, uintptr(unsafe.Pointer(&ver))); err != nil {
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

	if err := ioctl(c, ioctlVersion, uintptr(unsafe.Pointer(&ver))); err != nil {
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

func cToGoString(b []byte) string {
	return string(bytes.Trim(b, "\u0000"))
}
