package drm

import (
	"bytes"
	"fmt"
	"unsafe"
)

func (c *Card) ModeGetResources() (*ModeResources, error) {
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

func (c *Card) ModeGetCRTC(crtcID uint32) (*ModeCRTC, error) {
	crtc := cModeCRTC{CRTCID: crtcID}
	if err := ioctl(c.fd, ioctlModeGetCRTC, uintptr(unsafe.Pointer(&crtc))); err != nil {
		return nil, fmt.Errorf("ioctl: %w", err)
	}
	return &ModeCRTC{
		cModeCRTC: crtc,
		Name:      cToGoString(crtc.name[:]),
	}, nil
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

func (c *Card) ModeSetProperty(connectorID, propID uint32, value uint64) error {
	prop := cModeConnectorSetProperty{
		value:       value,
		propID:      propID,
		connectorID: connectorID,
	}
	if err := ioctl(c.fd, ioctlModeSetProperty, uintptr(unsafe.Pointer(&prop))); err != nil {
		return fmt.Errorf("ioctl: %w", err)
	}
	return nil
}
