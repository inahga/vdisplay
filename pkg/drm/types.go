package drm

// Most comments here are taken directly from drm/drm.h or drm/drm_mode.h

const (
	// ModePropPending is deprecated, do not use.
	ModePropPending uint32 = 1 << iota
	ModePropRange
	ModePropImmutable
	// ModePropEnum indicates that the type is enumerated with text strings.
	ModePropEnum
	ModePropBlob
	// ModePropBitmask indicates the type is a bitmask of enumerated types.
	ModePropBitmask
)

const (
	// ModePropLegacyType enumerates non-extended types: legacy bitmask, one bit
	// per type
	ModePropLegacyType = ModePropRange | ModePropEnum | ModePropBlob | ModePropBitmask

	// ModePropExtendedType enumerates extended types. Rather than continue to
	// consume a bit per type, grab a chunk of the bits to use as integer type id.
	ModePropExtendedType = 0x0000FFC0
	ModePropObject       = uint32(1 << 6)
	ModePropSignedRange  = uint32(2 << 6)

	// ModePropAtomic is used to hide properties from userspace that are not aware
	// of atomic properties. This is mostly to work around older userspace (DDX
	// drivers that read/write each prop they find without being aware this could
	// be triggering a lengthy modeset
	ModePropAtomic = 0x80000000
)

const (
	// ClientCapStereo3D indicates the DRM core will expose the stereo 3D capabilties
	// of the monitor by advertising the supported 3D layouts on the flags of struct
	// drm_mode_modeinfo, if set to 1.
	ClientCapStereo3D uint64 = iota + 1
	// ClientCapUniversalPlanes indicates the DRM core will expose all planes
	// (overlay, primary, and cursor) to userspace, if set to 1.
	ClientCapUniversalPlanes
	// ClientCapAtomic indicates the DRM core will expose atomic properties to
	// userspace, if set to 1.
	ClientCapAtomic
	// ClientCapAspectRatio indicates the DRM core will provide aspect ratio
	// information in modes, if set to 1.
	ClientCapAspectRatio
	// ClientCapWritebackConnectors indicates the DRM core will expose special
	// connectors to be used for writing back to memory the scene setup in the
	// commit, if set to 1. Depends on the client also supporting ClientCapAtomic.
	ClientCapWritebackConnectors
)

// This is for connectors with multiple signal types. Try to match ModeConnectorX
// as closely as possible.
const (
	ModeSubconnectorAutomatic uint32 = 0
	ModeSubconnectorUnknown          = 0
	ModeSubconnectorDVID             = 3
	ModeSubconnectorDVIA             = 4
	ModeSubconnectorComposite        = 5
	ModeSubconnectorSVIDEO           = 6
	ModeSubconnectorComponent        = 8
	ModeSubconnectorSCART            = 9
)

const (
	ModeConnectorUnknown uint32 = iota
	ModeConnectorVGA
	ModeConnectorDVII
	ModeConnectorDVID
	ModeConnectorDVIA
	ModeConnectorComposite
	ModeConnectorSVIDEO
	ModeConnectorLVDS
	ModeConnectorComponent
	ModeConnector9PinDIN
	ModeConnectorDisplayPort
	ModeConnectorHDMIA
	ModeConnectorHDMIB
	ModeConnectorTV
	ModeConnectorEDP
	ModeConnectorVirtual
	ModeConnectorDSI
	ModeConnectorDPI
	ModeConnectorWriteback
	ModeConnectorSPI
)

const (
	ModeEncoderNone uint32 = iota
	ModeEncoderDAC
	ModeEncoderTMDS
	ModeEncoderLVDS
	ModeEncoderTVDAC
	ModeEncoderVirtual
	ModeEncoderDSI
	ModeEncoderDPMST
	ModeEncoderDPI
)

const (
	ModeObjectCrtc      uint32 = 0xcccccccc
	ModeObjectConnector uint32 = 0xc0c0c0c0
	ModeObjectEncoder   uint32 = 0xe0e0e0e0
	ModeObjectMode      uint32 = 0xdededede
	ModeObjectProperty  uint32 = 0xb0b0b0b0
	ModeObjectFb        uint32 = 0xfbfbfbfb
	ModeObjectBlob      uint32 = 0xbbbbbbbb
	ModeObjectPlane     uint32 = 0xeeeeeeee
	ModeObjectAny       uint32 = 0
)

type Version struct {
	Major      int32
	Minor      int32
	PatchLevel int32
	Name       string
	Date       string
	Desc       string
}

type ModeResources struct {
	FBIDs        []uint32
	CRTCIDs      []uint32
	ConnectorIDs []uint32
	EncoderIDs   []uint32

	MinWidth  uint32
	MaxWidth  uint32
	MinHeight uint32
	MaxHeight uint32
}

type ModeInfo struct {
	cModeInfo
	Name string
}

type ModeConnector struct {
	cModeGetConnector

	EncoderIDs []uint32
	Modes      []ModeInfo
	PropIDs    []uint32
	PropValues []uint64
}

type ModeEncoder struct {
	cModeGetEncoder
}

type ModeCRTC struct {
	cModeCRTC
	Name string
	// SetConnectors is a list of connector IDs to be added when calling ModeSetCRTC()
	SetConnectors []uint32
}

type ModePropertyEnum struct {
	Value uint64
	Name  string
}

type ModeProperty struct {
	Values []uint64
	Enums  []ModePropertyEnum

	Name   string
	Flags  uint32
	PropID uint32
}

type ModeObjProperties struct {
	PropIDs    []uint32
	PropValues []uint64
	ID         uint32
	Type       uint32
}

type ModeBlob struct {
	ID   uint32
	Data []uint8
}

// ModePlaneResources is a slice of Plane IDs.
type ModePlaneResources []uint32

type ModePlane struct {
	cModeGetPlane
	FormatTypes []uint32
}

type ModeLease struct {
	Fd      uint32
	ID      uint32
	Objects []uint32
}

type ModeDumbBuffer struct {
	cModeCreateDumb
}

type ModeFramebuffer struct {
	cModeFBCmd
}
