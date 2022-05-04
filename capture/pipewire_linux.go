package capture

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/godbus/dbus/v5"
)

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

const (
	dbusSourceTypeMonitor uint32 = 1 << iota
	dbusSourceTypeWindow
	dbusSourceTypeVirtual
)

func NewPipewire() (ret *PipewireStream, err error) {
	ret = &PipewireStream{}
	ret.dbusConn, err = dbus.ConnectSessionBus()
	if err != nil {
		return nil, fmt.Errorf("pipewire: dbus: %w", err)
	}
	log.Printf("[pipewire] opened dbus connection: %s", ret.dbusConn.Names()[0])
	return ret, nil
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
	log.Printf("[pipewire] pipewire id is %d", p.streams[0].NodeID)
	// get pipewire remote fd, that's how we'll connect to it
	// need to gracefully handle cancellation of session
	return nil
}

func (p *PipewireStream) Close() error {
	return p.dbusConn.Close()
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
		args: []interface{}{vardict{
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
		args: []interface{}{
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
		args: []interface{}{
			p.sessionHandle,
			"", // TODO: select the parent window of the GUI app
			vardict{"handle_token": dbus.MakeVariant(handleToken)},
		},
	})
}

type dbusRequest struct {
	dest, method, handleToken string
	path                      dbus.ObjectPath
	flags                     dbus.Flags
	processResponse           func(vardict) error
	args                      []interface{}
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

func genToken(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:]
}

// See https://flatpak.github.io/xdg-desktop-portal/#gdbus-org.freedesktop.portal.Request
func uniqueNameToPath(name string) string {
	return strings.ReplaceAll(name[1:], ".", "_")
}
