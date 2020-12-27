package drm

import (
	"os"
	"syscall"
	"unsafe"
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

type unimplemented struct{}

var (
	ioctlVersion      = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cVersion{})), ioctlBase, 0x00)
	ioctlGetUnique    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x01)
	ioctlGetMagic     = ioctlRequest(iocRead, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x02)
	ioctlIrqBusid     = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x03)
	ioctlGetMap       = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x04)
	ioctlGetClient    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x05)
	ioctlGetStats     = ioctlRequest(iocRead, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x06)
	ioctlSetVersion   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x07)
	ioctlModesetCtl   = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x08)
	ioctlGemClose     = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x09)
	ioctlGemFlink     = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x0a)
	ioctlGemOpen      = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x0b)
	ioctlGetCap       = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x0c)
	ioctlSetClientCap = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(cSetClientCap{})), ioctlBase, 0x0d)

	ioctlSetUnique = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x10)
	ioctlAuthMagic = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x11)
	ioctlBlock     = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x12)
	ioctlUnblock   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x13)
	ioctlControl   = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x14)
	ioctlAddMap    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x15)
	ioctlAddBufs   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x16)
	ioctlMarkBufs  = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x17)
	ioctlInfoBufs  = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x18)
	ioctlMapBufs   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x19)
	ioctlFreeBufs  = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x1a)

	ioctlRmMap = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x1b)

	ioctlSetSAreaCtx = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x1c)
	ioctlGetSAreaCtx = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x1d)

	ioctlSetMaster  = ioctlRequest(iocNone, 0, ioctlBase, 0x1e)
	ioctlDropMaster = ioctlRequest(iocNone, 0, ioctlBase, 0x1f)

	ioctlAddCtx    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x20)
	ioctlRmCtx     = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x21)
	ioctlModCtx    = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x22)
	ioctlGetCtx    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x23)
	ioctlSwitchCtx = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x24)
	ioctlNewCtx    = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x25)
	ioctlResCtx    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x26)
	ioctlAddDraw   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x27)
	ioctlRmDraw    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x28)
	ioctlDMA       = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x29)
	ioctlLock      = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x2a)
	ioctlUnlock    = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x2b)
	ioctlFinish    = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x2c)

	ioctlPrimeHandleToFd = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x2d)
	ioctlPrimeFdToHandle = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x2e)

	ioctlAGPAcquire = ioctlRequest(iocNone, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x30)
	ioctlAGPRelease = ioctlRequest(iocNone, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x31)
	ioctlAGPEnable  = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x32)
	ioctlAGPInfo    = ioctlRequest(iocRead, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x33)
	ioctlAGPAlloc   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x34)
	ioctlAGPFree    = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x35)
	ioctlAGPBind    = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x36)
	ioctlAGPUnbind  = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x37)

	ioctlSGAlloc = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x38)
	ioctlSGFree  = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x39)

	ioctlWaitVblank = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x3a)

	ioctlCRTCGetSequence   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x3b)
	ioctlCRTCQueueSequence = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x3c)

	ioctlUpdateDraw = ioctlRequest(iocWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0x3f)

	ioctlModeGetResources = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeCardRes{})), ioctlBase, 0xA0)
	ioctlModeGetCRTC      = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeCRTC{})), ioctlBase, 0xA1)
	ioctlModeSetCRTC      = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeCRTC{})), ioctlBase, 0xA2)
	ioctlModeCursor       = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xA3)
	ioctlModeGetGamma     = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xA4)
	ioctlModeSetGamma     = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xA5)
	ioctlModeGetEncoder   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeGetEncoder{})), ioctlBase, 0xA6)
	ioctlModeGetConnector = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeGetConnector{})), ioctlBase, 0xA7)
	ioctlModeAttachMode   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xA8)
	ioctlModeDetachMode   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xA9)

	ioctlModeGetProperty = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeGetProperty{})), ioctlBase, 0xAA)
	ioctlModeSetProperty = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeConnectorSetProperty{})), ioctlBase, 0xAB)
	ioctlModeGetPropBlob = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeGetBlob{})), ioctlBase, 0xAC)
	ioctlModeGetFB       = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeFBCmd{})), ioctlBase, 0xAD)
	ioctlModeAddFB       = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeFBCmd{})), ioctlBase, 0xAE)
	ioctlModeRmFB        = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(uint32(0))), ioctlBase, 0xAF)
	ioctlModePageFlip    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xB0)
	ioctlModeDirtyFB     = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xB1)

	ioctlModeCreateDumb        = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeCreateDumb{})), ioctlBase, 0xB2)
	ioctlModeMapDumb           = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeMapDumb{})), ioctlBase, 0xB3)
	ioctlModeDestroyDumb       = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeDestroyDumb{})), ioctlBase, 0xB4)
	ioctlModeGetPlaneResources = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeGetPlaneRes{})), ioctlBase, 0xB5)
	ioctlModeGetPlane          = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeGetPlane{})), ioctlBase, 0xB6)
	ioctlModeSetPlane          = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xB7)
	ioctlModeAddFB2            = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xB8)
	ioctlModeObjGetProperties  = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeObjGetProperties{})), ioctlBase, 0xB9)
	ioctlModeObjSetProperty    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xBA)
	ioctlModeCursor2           = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xBB)
	ioctlModeAtomic            = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xBC)
	ioctlModeCreatePropBlob    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xBD)
	ioctlModeDestroyPropBlob   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xBE)

	ioctlSyncObjCreate     = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xBF)
	ioctlSyncObjDestroy    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xC0)
	ioctlSyncObjHandleToFd = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xC1)
	ioctlSyncObjFdToHandle = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xC2)
	ioctlSyncObjWait       = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xC3)
	ioctlSyncObjReset      = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xC4)
	ioctlSyncObjSignal     = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xC5)

	ioctlModeCreateLease = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeCreateLease{})), ioctlBase, 0xC6)
	ioctlModeListLessees = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeListLessees{})), ioctlBase, 0xC7)
	ioctlModeGetLease    = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeGetLease{})), ioctlBase, 0xC8)
	ioctlModeRevokeLease = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(cModeRevokeLease{})), ioctlBase, 0xC9)

	ioctlSyncObjTimelineWait   = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xCA)
	ioctlSyncObjQuery          = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xCB)
	ioctlSyncObjTransfer       = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xCC)
	ioctlSyncObjTimelineSignal = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xCD)

	ioctlModeGetFB2 = ioctlRequest(iocReadWrite, uint16(unsafe.Sizeof(unimplemented{})), ioctlBase, 0xCE)
)
