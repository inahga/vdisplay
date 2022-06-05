package drm

import (
	"os"
	"syscall"
)

const (
	// IoctlBase is the DRM specific character that identifies DRM ioctls.
	ioctlBase uint8 = 'd'

	iocNone = 0
	// iocWrite indicates that userspace is writing, and kernel is reading.
	iocWrite = 0x1
	// iocRead indicates that userland is reading, and kernel is writing.
	iocRead      = 0x2
	iocReadWrite = iocRead | iocWrite

	iocNRBits   = 8
	iocTypeBits = 8
	iocSizeBits = 14
	iocDirBits  = 2

	iocNRMask   = (1 << iocNRBits) - 1
	iocTypeMask = (1 << iocTypeBits) - 1
	iocSizeMask = (1 << iocSizeBits) - 1
	iocDirMask  = (1 << iocDirBits) - 1

	iocNRShift   = 0
	iocTypeShift = iocNRShift + iocNRBits
	iocSizeShift = iocTypeShift + iocTypeBits
	iocDirShift  = iocSizeShift + iocSizeBits
)

// ioctlRequest generates an ioctl request magic number. Dir corresponds to the
// read/write direction. Size corresponds to the corresponding ioctl data structure
// size, Typ corresponds to the character that identifies the driver (i.e. 'd' for
// DRM). Nr corresponds to the function magic number.
func ioctlRequest(dir uint8, size uint16, typ, nr uint8) uint32 {
	return (uint32(dir) << iocDirShift) | (uint32(size) << iocSizeShift) |
		(uint32(typ) << iocTypeShift) | (uint32(nr) << iocNRShift)
}

func ioctl(fd *os.File, request uint32, data uintptr) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(), uintptr(request), data); err != 0 {
		return err
	}
	return nil
}
