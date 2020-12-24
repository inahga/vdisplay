// Package drm provides functions and structs for interacting with the Linux
// kernel's Direct Rendering Manager.
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
		ver.name = uintptr(unsafe.Pointer(&name[0]))
	}
	if ver.datelen > 0 {
		date = make([]byte, ver.datelen+1)
		ver.date = uintptr(unsafe.Pointer(&date[0]))
	}
	if ver.desclen > 0 {
		desc = make([]byte, ver.desclen+1)
		ver.desc = uintptr(unsafe.Pointer(&desc[0]))
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

func (c *Card) ModeResources() (*ModeResources, error) {
	var res cModeCardRes
	if err := ioctl(c.fd, ioctlModeGetResources, uintptr(unsafe.Pointer(&res))); err != nil {
		return nil, fmt.Errorf("ioctl: %w", err)
	}

	ret := ModeResources{
		MinWidth:  res.minWidth,
		MaxWidth:  res.maxWidth,
		MinHeight: res.minHeight,
		MaxHeight: res.maxHeight,
	}
	if res.countConnectors > 0 {
		ret.ConnectorIDs = make([]uint32, res.countConnectors)
		res.connectorIDPtr = uintptr(unsafe.Pointer(&ret.ConnectorIDs[0]))
	}
	if res.countCRTC > 0 {
		ret.CRTCIDs = make([]uint32, res.countCRTC)
		res.crtcIDPtr = uintptr(unsafe.Pointer(&ret.CRTCIDs[0]))
	}
	if res.countEncoders > 0 {
		ret.EncoderIDs = make([]uint32, res.countEncoders)
		res.encoderIDPtr = uintptr(unsafe.Pointer(&ret.EncoderIDs[0]))
	}
	if res.countFB > 0 {
		ret.FBIDs = make([]uint32, res.fbIDPtr)
		res.fbIDPtr = uintptr(unsafe.Pointer(&ret.FBIDs[0]))
	}
	if err := ioctl(c.fd, ioctlModeGetResources, uintptr(unsafe.Pointer(&res))); err != nil {
		return nil, fmt.Errorf("ioctl: %w", err)
	}

	return &ret, nil
}

func (c *Card) ModeGetConnector(connectorID uint32) (*ModeConnector, error) {
	conn := cModeGetConnector{ConnectorID: connectorID}
	if err := ioctl(c.fd, ioctlModeGetConnector, uintptr(unsafe.Pointer(&conn))); err != nil {
		return nil, fmt.Errorf("ioctl: %w", err)
	}

	var modes []cModeInfo
	ret := ModeConnector{
		cModeGetConnector: conn,
	}
	if conn.countEncoders > 0 {
		ret.EncoderIDs = make([]uint32, conn.countEncoders)
		conn.encodersPtr = uintptr(unsafe.Pointer(&ret.EncoderIDs[0]))
	}
	if conn.countModes > 0 {
		modes = make([]cModeInfo, conn.countModes)
		conn.modesPtr = uintptr(unsafe.Pointer(&modes[0]))
	}
	if conn.countProps > 0 {
		ret.PropIDs = make([]uint32, conn.countProps)
		ret.PropValues = make([]uint64, conn.countProps)
		conn.propsPtr = uintptr(unsafe.Pointer(&ret.PropIDs[0]))
		conn.propValuesPtr = uintptr(unsafe.Pointer(&ret.PropValues[0]))
	}
	if err := ioctl(c.fd, ioctlModeGetConnector, uintptr(unsafe.Pointer(&conn))); err != nil {
		return nil, fmt.Errorf("ioctl: %w", err)
	}

	for _, mode := range modes {
		ret.Modes = append(ret.Modes, ModeInfo{
			cModeInfo: mode,
			Name:      cToGoString(mode.name[:]),
		})
	}
	return &ret, nil
}

func (c *Card) ModeGetProperty(propID uint32) (*ModeProperty, error) {
	prop := cModeGetProperty{propID: propID}
	if err := ioctl(c.fd, ioctlModeGetProperty, uintptr(unsafe.Pointer(&prop))); err != nil {
		return nil, fmt.Errorf("ioctl: %w", err)
	}

	var enums []cModePropertyEnum
	ret := ModeProperty{
		PropID: prop.propID,
		Flags:  prop.flags,
		Name:   string(bytes.Trim(prop.name[:], "\u0000")),
	}
	if prop.countValues > 0 {
		ret.Values = make([]uint64, prop.countValues)
		prop.valuesPtr = uintptr(unsafe.Pointer(&ret.Values[0]))
	}
	if prop.countEnumBlobs > 0 {
		enums = make([]cModePropertyEnum, prop.countEnumBlobs)
		prop.enumBlobPtr = uintptr(unsafe.Pointer(&enums[0]))
	}
	if err := ioctl(c.fd, ioctlModeGetProperty, uintptr(unsafe.Pointer(&prop))); err != nil {
		return nil, fmt.Errorf("ioctl: %w", err)
	}

	for _, enum := range enums {
		ret.Enums = append(ret.Enums, ModePropertyEnum{
			Value: enum.value,
			Name:  cToGoString(enum.name[:]),
		})
	}
	return &ret, nil
}
