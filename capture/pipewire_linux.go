package capture

// #cgo pkg-config: libpipewire-0.3
// #include <pipewire/pipewire.h>
//
// int pipewire_init(uint32_t, uint32_t, uint32_t);
import "C"
import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
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

	sendCh chan<- Buffer
	endCh  chan struct{}

	maxFramerate uint32
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
	pipewireReceiverMap = map[uint32]chan<- Buffer{}
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

func (p *PipewireStream) Register(ch chan<- Buffer) {
	// TODO: handle the case where a channel is already registered in the global map
	p.sendCh = ch
}

func (p *PipewireStream) SetMaxFramerate(f uint32) {
	p.maxFramerate = f
}

func (p *PipewireStream) Start() (err error) {
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

	// TODO: we probably need to lock this, or make it so that you can't re-assign
	// the channel
	pipewireReceiverMap[p.streams[0].NodeID] = p.sendCh

	// For now we are only concerned with the first stream node ID.
	// Can have multiple, but we did not set that up in selectSources()
	go func() {
		ret := C.pipewire_init(C.uint(p.streamFD), C.uint(p.streams[0].NodeID), C.uint(p.maxFramerate))
		if ret < 0 {
			// TODO: cleaner error handling here
			panic(fmt.Errorf("[pipewire] pipewire_init exit status %d", ret))
		}
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

//export receiveBuffer
func receiveBuffer(nodeID C.uint, b *C.struct_pw_buffer) {
	_, ok := pipewireReceiverMap[uint32(nodeID)]
	if !ok {
		log.Printf("[pipewire] received buffer for unknown channel for pipewire node ID %d", nodeID)
		return
	}
	log.Printf("[pipewire] received buffer into go for ID %d", nodeID)
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
