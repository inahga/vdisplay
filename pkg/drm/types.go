package drm

const (
	// ModePropPending is deprecated, do not use.
	ModePropPending uint32 = 1 << iota
	ModePropRange
	ModePropImmutable
	// ModePropEnum indicates that the type is enumerated with text strings.
	ModePropEnum
	ModePropBlob
	// ModePropBitmask indicates the type is a bitmask of enumerated types (?)
	ModePropBitmask
)

// Comments here are copied from drm/drm_mode.h
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

type Version struct {
	Major      int32
	Minor      int32
	PatchLevel int32
	Name       string
	Date       string
	Desc       string
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
