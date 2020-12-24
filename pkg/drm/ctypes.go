package drm

type (
	kernelSize = uint64
	cint       = int32
)

const (
	connectorNameLen = 32
	displayModeLen   = 32
	propNameLen      = 32
)

type cVersion struct {
	major      cint
	minor      cint
	patchlevel cint
	namelen    kernelSize
	name       uintptr // to a []byte
	datelen    kernelSize
	date       uintptr // to a []byte
	desclen    kernelSize
	desc       uintptr // to a []byte
}

type cSetClientCap struct {
	capability uint64
	value      uint64
}

type cModeInfo struct {
	Clock      uint32
	HDisplay   uint16
	HSyncStart uint16
	HSyncEnd   uint16
	HTotal     uint16
	HSkew      uint16
	VDisplay   uint16
	VSyncStart uint16
	VSyncEnd   uint16
	VTotal     uint16
	VScan      uint16
	VRefresh   uint32
	Flags      uint32
	Type       uint32
	name       [displayModeLen]byte
}

type cModeGetConnector struct {
	encodersPtr   uintptr // to a []uint32
	modesPtr      uintptr // to a []cModeInfo
	propsPtr      uintptr // to a []uint32
	propValuesPtr uintptr // to a []uint64
	countModes    uint32
	countProps    uint32
	countEncoders uint32

	EncoderID       uint32 // current encoder (so says drm/drm_mode.h)
	ConnectorID     uint32
	ConnectorType   uint32
	ConnectorTypeID uint32
	Connection      uint32
	MMWidth         uint32
	MMHeight        uint32
	Subpixel        uint32
	Pad             uint32
}

type cModeCardRes struct {
	fbIDPtr         uintptr // to a []uint32
	crtcIDPtr       uintptr // to a []uint32
	connectorIDPtr  uintptr // to a []uint32
	encoderIDPtr    uintptr // to a []uint32
	countFB         uint32
	countCRTC       uint32
	countConnectors uint32
	countEncoders   uint32

	minWidth  uint32
	maxWidth  uint32
	minHeight uint32
	maxHeight uint32
}

type cModePropertyEnum struct {
	value uint64
	name  [propNameLen]byte
}

type cModeGetProperty struct {
	valuesPtr   uintptr // Values and blob lengths
	enumBlobPtr uintptr // Enum and blob ID ptrs

	propID uint32
	flags  uint32
	name   [propNameLen]byte

	countValues uint32
	// This is only used to count enum values, not blobs. The blobs is simply
	// because of historical reason, i.e. backwards compat.
	countEnumBlobs uint32
}

type cModeConnectorSetProperty struct {
	value       uint64
	propID      uint32
	connectorID uint32
}

type cModeCRTC struct {
	setConnectorsPtr uintptr // to a []uint32
	countConnectors  uint32

	CRTCID    uint32
	FBID      uint32
	X         uint32 // x position on the framebuffer
	Y         uint32 // y position on the framebuffer
	GammaSize uint32
	ModeValid uint32
	cModeInfo
}

type cModeObjGetProperties struct {
	propsPtr      uintptr // to a []uint32
	propValuesPtr uintptr // to a []uint64
	countProps    uint32
	objID         uint32
	objType       uint32
}
