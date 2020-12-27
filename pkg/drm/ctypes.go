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
	name       uint64 // ptr to a []byte
	datelen    kernelSize
	date       uint64 // ptr to a []byte
	desclen    kernelSize
	desc       uint64 // ptr to a []byte
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
	encodersPtr   uint64 // ptr to a []uint32
	modesPtr      uint64 // ptr to a []cModeInfo
	propsPtr      uint64 // ptr to a []uint32
	propValuesPtr uint64 // ptr to a []uint64
	countModes    uint32
	countProps    uint32
	countEncoders uint32

	EncoderID  uint32 // current encoder (so says drm/drm_mode.h)
	ID         uint32
	Type       uint32
	TypeID     uint32
	Connection uint32
	MMWidth    uint32
	MMHeight   uint32
	Subpixel   uint32
	pad        uint32
}

type cModeCardRes struct {
	fbIDPtr         uint64 // ptr to a []uint32
	crtcIDPtr       uint64 // ptr to a []uint32
	connectorIDPtr  uint64 // ptr to a []uint32
	encoderIDPtr    uint64 // ptr to a []uint32
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
	valuesPtr   uint64 // Values and blob lengths
	enumBlobPtr uint64 // Enum and blob ID ptrs

	propID uint32
	flags  uint32
	name   [propNameLen]byte

	countValues uint32
	// This is only used ptr to count enum values, not blobs. The blobs is simply
	// because of historical reason, i.e. backwards compat.
	countEnumBlobs uint32
}

type cModeConnectorSetProperty struct {
	value       uint64
	propID      uint32
	connectorID uint32
}

type cModeCRTC struct {
	setConnectorsPtr uint64 // ptr to a []uint32
	countConnectors  uint32

	ID        uint32
	FBID      uint32
	X         uint32 // x position on the framebuffer
	Y         uint32 // y position on the framebuffer
	GammaSize uint32
	ModeValid uint32
	cModeInfo
}

type cModeGetEncoder struct {
	ID             uint32
	Type           uint32
	CRTCID         uint32
	PossibleCRTCs  uint32
	PossibleClones uint32
}

type cModeObjGetProperties struct {
	propsPtr      uint64 // ptr to a []uint32
	propValuesPtr uint64 // ptr to a []uint64
	countProps    uint32
	objID         uint32
	objType       uint32
}

type cModeGetBlob struct {
	blobID uint32
	length uint32
	data   uint64 // ptr to a []uint8
}

type cModeGetPlaneRes struct {
	planeIDPtr  uint64
	countPlanes uint32 // ptr to a []uint32
}

type cModeGetPlane struct {
	ID            uint32
	CRTCID        uint32
	FBID          uint32
	PossibleCRTCs uint32
	GammaSize     uint32

	countFormatTypes uint32
	formatTypePtr    uint64 // ptr to a []uint32
}

type cModeCreateLease struct {
	objectIDs   uint64 // ptr to a []uint32
	objectCount uint32
	flags       uint32 // flags for new file descriptor

	lesseeID uint32
	fd       uint32
}

type cModeGetLease struct {
	countObjects uint32
	pad          uint32
	objectsPtr   uint64 // ptr to a []uint32
}

type cModeListLessees struct {
	countLessees uint32
	pad          uint32
	lesseesPtr   uint64 // ptr to a []uint32
}

type cModeRevokeLease struct {
	lesseeID uint32
}

type cModeCreateDumb struct {
	Height uint32
	Width  uint32
	Bpp    uint32
	flags  uint32 // unused

	Handle uint32
	Pitch  uint32
	Size   uint64
}

type cModeMapDumb struct {
	handle uint32
	pad    uint32
	offset uint64 // fake offset to use for subsequent mmap call
}

type cModeDestroyDumb struct {
	handle uint32
}

type cModeFBCmd struct {
	ID     uint32
	Width  uint32
	Height uint32
	Pitch  uint32
	Bpp    uint32
	Depth  uint32
	Handle uint32 // driver specific handle to a buffer
}
