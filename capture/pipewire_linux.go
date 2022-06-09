package capture

/*
#cgo pkg-config: libpipewire-0.3
#include <pipewire/pipewire.h>
#include <spa/buffer/buffer.h>
#include <spa/param/video/format-utils.h>
#include <spa/param/video/type-info.h>

int pipewire_run_loop(uint32_t, uint32_t, uint32_t);
*/
import "C"
import (
	"crypto/rand"
	"errors"
	"fmt"
	"image"
	"log"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/godbus/dbus/v5"
	"github.com/inahga/vdisplay/internal/convert"
)

// OBS Studio source code is a good reference as to how this is done.
// https://github.com/obsproject/obs-studio/blob/3d7663f417d92d20576ccc7fe455d11e25ebf5a9/plugins/linux-pipewire/pipewire.c

// PipewireStream uses the org.freedesktop.portal.ScreenCast portal granted by
// xdg-desktop-portal to generate a pipewire stream for screen capture.
type PipewireStream struct {
	dbusConn      *dbus.Conn
	sessionHandle dbus.ObjectPath
	restoreToken  string
	streams       []struct {
		NodeID     uint32
		Properties vardict
	}
	streamFD dbus.UnixFDIndex

	endCh chan struct{}
	cb    func(image.Image)
}

type vardict = map[string]dbus.Variant

var (
	ErrDbusBadResponse   = errors.New("unknown or malformed response")
	ErrDbusUserCancelled = errors.New("user cancelled interaction")
	ErrDbusCancelled     = errors.New("interaction cancelled")
)

const (
	dbusCursorModeHidden uint32 = 1 << iota
	dbusCursorModeEmbedded
	dbusCursorModeMetadata
)

var (
	// Because we can't pass go methods of complex structs to cgo, we will identify
	// the channels by their pipewire ID. Only one channel is supported.
	pipewireReceiverMap = map[uint32]*PipewireStream{}
	// pipewireReceiverMapLock sync.Mutex
)

func init() {
	log.Println("[pipewire] pw_init()")
	C.pw_init(nil, nil)
}

func NewPipewire() (ret *PipewireStream, err error) {
	ret = &PipewireStream{endCh: make(chan struct{})}
	ret.dbusConn, err = dbus.ConnectSessionBus()
	if err != nil {
		return nil, fmt.Errorf("pipewire: dbus: %w", err)
	}
	log.Printf("[pipewire] opened dbus connection: %s", ret.dbusConn.Names()[0])
	return ret, nil
}

func (p *PipewireStream) Close() error {
	// TODO: more needs to be done here
	// close pipewire thread and stream
	// close dbus session
	// unregister pipewire stream from global channel map
	return p.dbusConn.Close()
}

func (p *PipewireStream) Start(framerate uint32, _ image.Rectangle, cb func(image.Image)) (err error) {
	if err := p.createSession(); err != nil {
		return fmt.Errorf("createSession: %w", err)
	}
	if err := p.selectSources(); err != nil {
		return fmt.Errorf("selectSources: %w", err)
	}
	log.Printf("[pipewire] created dbus screencast session")
	if err := p.startSession(); err != nil {
		return fmt.Errorf("startSession: %w", err)
	}
	log.Printf("[pipewire] cast id is %d", p.streams[0].NodeID)
	if err := p.getStreamFD(); err != nil {
		return fmt.Errorf("getStreamFD: %w", err)
	}
	log.Printf("[pipewire] cast fd is %d", p.streamFD)

	// TODO: we probably need to lock this
	pipewireReceiverMap[p.streams[0].NodeID] = p

	// For now we are only concerned with the first stream node ID.
	// Can have multiple, but we did not set that up in selectSources()
	p.cb = cb
	go func() {
		ret := C.pipewire_run_loop(C.uint(p.streamFD), C.uint(p.streams[0].NodeID), C.uint(framerate))
		// todo: cleaner exit handling
		panic(fmt.Errorf("[pipewire] pipewire_init exit status %d", ret))
	}()

	<-time.After(time.Second * 3)
	// TODO: create a listener for dbus stream close events
	return nil
}

func (p *PipewireStream) createSession() error {
	sessionHandleToken, handleToken := genToken(16), genToken(16)
	return p.dbusRequest(&dbusRequest{
		dest:        "org.freedesktop.portal.Desktop",
		path:        "/org/freedesktop/portal/desktop",
		method:      "org.freedesktop.portal.ScreenCast.CreateSession",
		flags:       0,
		handleToken: handleToken,
		processResponse: func(vardict vardict) error {
			if handle, ok := vardict["session_handle"]; ok {
				handle.Store(&p.sessionHandle)
			} else {
				return fmt.Errorf("%w: missing session_handle", ErrDbusBadResponse)
			}
			return nil
		},
		args: []any{vardict{
			"handle_token":         dbus.MakeVariant(handleToken),
			"session_handle_token": dbus.MakeVariant(sessionHandleToken),
		}},
	})
}

func (p *PipewireStream) selectSources() error {
	handleToken := genToken(16)
	return p.dbusRequest(&dbusRequest{
		dest:        "org.freedesktop.portal.Desktop",
		path:        "/org/freedesktop/portal/desktop",
		method:      "org.freedesktop.portal.ScreenCast.SelectSources",
		flags:       0,
		handleToken: handleToken,
		args: []any{
			p.sessionHandle,
			vardict{
				"handle_token": dbus.MakeVariant(handleToken),
				// TODO: need to check if cursor mode is available first
				"cursor_mode": dbus.MakeVariant(dbusCursorModeEmbedded),
				// TODO: session persistence
			},
		},
	})
}

func (p *PipewireStream) startSession() error {
	handleToken := genToken(16)
	return p.dbusRequest(&dbusRequest{
		dest:        "org.freedesktop.portal.Desktop",
		path:        "/org/freedesktop/portal/desktop",
		method:      "org.freedesktop.portal.ScreenCast.Start",
		flags:       0,
		handleToken: handleToken,
		processResponse: func(vardict vardict) error {
			if handle, ok := vardict["streams"]; ok {
				handle.Store(&p.streams)
			} else {
				return fmt.Errorf("%w: missing streams", ErrDbusBadResponse)
			}
			if handle, ok := vardict["restore_token"]; ok {
				handle.Store(&p.restoreToken)
			}
			return nil
		},
		args: []any{
			p.sessionHandle,
			"", // TODO: select the parent window of the GUI app
			vardict{"handle_token": dbus.MakeVariant(handleToken)},
		},
	})
}

func (p *PipewireStream) getStreamFD() error {
	return p.dbusConn.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop").
		Call("org.freedesktop.portal.ScreenCast.OpenPipeWireRemote", 0, p.sessionHandle, vardict{}).
		Store(&p.streamFD)
}

type dbusRequest struct {
	dest, method, handleToken string
	path                      dbus.ObjectPath
	flags                     dbus.Flags
	processResponse           func(vardict) error
	args                      []any
}

func (p *PipewireStream) dbusRequest(request *dbusRequest) error {
	matchRequestSignal := []dbus.MatchOption{
		dbus.WithMatchObjectPath(dbus.ObjectPath(string(request.path) + "/request/" +
			uniqueNameToPath(p.dbusConn.Names()[0]) + "/" + request.handleToken)),
		dbus.WithMatchInterface("org.freedesktop.portal.Request.Response"),
	}
	if err := p.dbusConn.AddMatchSignal(matchRequestSignal...); err != nil {
		return err
	}
	defer func() {
		if rerr := p.dbusConn.RemoveMatchSignal(matchRequestSignal...); rerr != nil {
			panic(rerr)
		}
	}()
	signal := make(chan *dbus.Signal)
	p.dbusConn.Signal(signal)
	defer func() {
		p.dbusConn.RemoveSignal(signal)
		close(signal)
	}()

	var requestHandle dbus.ObjectPath
	if err := p.dbusConn.Object(request.dest, request.path).
		Call(request.method, request.flags, request.args...).Store(&requestHandle); err != nil {
		return fmt.Errorf("call: %w", err)
	}

	log.Printf("[pipewire] awaiting dbus response")
	response := <-signal
	return checkResponseSignal(response, request.processResponse)
}

func checkResponseSignal(signal *dbus.Signal, resultsFn func(vardict) error) error {
	if len(signal.Body) < 2 {
		return ErrDbusBadResponse
	}
	if code, ok := signal.Body[0].(uint32); ok {
		switch code {
		case 0:
		case 1:
			return ErrDbusUserCancelled
		case 2:
			return ErrDbusCancelled
		default:
			return fmt.Errorf("%w: unknown response code %d", ErrDbusBadResponse, code)
		}
	} else {
		return ErrDbusBadResponse
	}
	if resultsFn != nil {
		if vardict, ok := signal.Body[1].(vardict); ok {
			if err := resultsFn(vardict); err != nil {
				return err
			}
		} else {
			return ErrDbusBadResponse
		}
	}
	return nil
}

//export pipewire_receive_buffer
func pipewire_receive_buffer(nodeID C.uint, format *C.struct_spa_video_info, b *C.struct_pw_buffer) {
	stream, ok := pipewireReceiverMap[uint32(nodeID)]
	if !ok {
		panic(fmt.Errorf("[pipewire] received buffer for unknown channel for pipewire node ID %d", nodeID))
	}
	log.Printf("[pipewire] received buffer into go for ID %d", nodeID)

	var (
		robuf, rwbuf []uint8
		err          error
		data         = b.buffer.datas
		meta         = b.buffer.metas
	)

	if b.buffer.n_datas > 1 {
		log.Printf("[pipewire] unexpectedly receiving more data, n_datas = %d", b.buffer.n_datas)
	}
	if b.buffer.n_metas > 1 {
		log.Printf("[pipewire] unexpectedly receiving more meta, n_metas = %d", b.buffer.n_metas)
	}
	if data.flags&C.SPA_DATA_FLAG_READABLE == 0 {
		panic(fmt.Errorf("buffer not readable, data flags = %d", data.flags))
	}
	if meta._type != C.SPA_META_Busy {
		log.Printf("[pipewire] unhandled meta type %d", meta._type)
	}

	switch data._type {
	case C.SPA_DATA_MemFd:
		robuf, err = syscall.Mmap(int(data.fd), int64(data.mapoffset), int(data.maxsize), syscall.PROT_READ, syscall.MAP_SHARED)
		if err != nil {
			panic(fmt.Errorf("mmap: %w", err))
		}
		log.Printf("[pipewire] mmap'd fd %d, size = %d", data.fd, data.maxsize)
		defer func() {
			if err := syscall.Munmap(robuf); err != nil {
				log.Printf("[pipewire] munmap: %s", err)
			}
		}()
	default:
		panic(fmt.Errorf("unsupported pipewire spa_data type %d", data._type))
	}

	rwbuf = make([]byte, len(robuf))
	copy(rwbuf, robuf)

	rawInfo := *(*C.struct_spa_video_info_raw)(unsafe.Pointer(&format.info[0]))
	switch rawInfo.format {
	case C.SPA_VIDEO_FORMAT_BGRx:
		convert.BGRxToRGBA(rwbuf)
	default:
		panic(fmt.Errorf("unsupported buffer format type %d", rawInfo.format))
	}

	img := image.NewRGBA(image.Rect(0, 0, int(rawInfo.size.width), int(rawInfo.size.height)))
	img.Pix = rwbuf
	stream.cb(img)
}

func genToken(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:]
}

// See https://flatpak.github.io/xdg-desktop-portal/#gdbus-org.freedesktop.portal.Request
func uniqueNameToPath(name string) string {
	return strings.ReplaceAll(name[1:], ".", "_")
}
