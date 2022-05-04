package capture

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
)

// Pipewire uses the org.freedesktop.portal.ScreenCast portal granted by
// xdg-desktop-portal to generate a pipewire stream for screen capture.
type Pipewire struct {
	sessionHandleToken string
	dbusConn           *dbus.Conn
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

func NewPipewire() (ret *Pipewire, err error) {
	ret = &Pipewire{}
	ret.dbusConn, err = dbus.ConnectSessionBus()
	if err != nil {
		return nil, fmt.Errorf("pipewire: dbus: %w", err)
	}
	return ret, nil
}

func checkResponseSignal(signal *dbus.Signal, resultsFn func(vardict)) error {
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
			resultsFn(vardict)
		} else {
			return ErrDbusBadResponse
		}
	}
	return nil
}

func (p *Pipewire) dbusRequest(dest string, path dbus.ObjectPath, method string, flags dbus.Flags, handleToken string, processResponse func(vardict), args ...interface{}) error {
	matchRequestSignal := []dbus.MatchOption{
		dbus.WithMatchObjectPath(dbus.ObjectPath(string(path) + "/request/" + uniqueNameToPath(p.dbusConn.Names()[0]) + "/" + handleToken)),
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
	if err := p.dbusConn.Object(dest, path).Call(method, 0, args...).Store(&requestHandle); err != nil {
		return fmt.Errorf("call: %w", err)
	}

	response := <-signal
	return checkResponseSignal(response, processResponse)
}

func (p *Pipewire) createSession() (sessionHandle dbus.ObjectPath, err error) {
	sessionHandleToken, handleToken := genToken(16), genToken(16)
	if err := p.dbusRequest(
		"org.freedesktop.portal.Desktop",
		"/org/freedesktop/portal/desktop",
		"org.freedesktop.portal.ScreenCast.CreateSession",
		0,
		handleToken,
		func(vardict vardict) {
			if handle, ok := vardict["session_handle"]; ok {
				handle.Store(&sessionHandle)
			}
		},
		vardict{
			"handle_token":         dbus.MakeVariant(handleToken),
			"session_handle_token": dbus.MakeVariant(sessionHandleToken),
		},
	); err != nil {
		return "", err
	}
	if sessionHandle == "" {
		return "", fmt.Errorf("%w: missing session_handle", ErrDbusBadResponse)
	}
	return sessionHandle, nil
}

func (p *Pipewire) selectSources(session dbus.ObjectPath) error {
	handleToken := genToken(16)
	return p.dbusRequest(
		"org.freedesktop.portal.Desktop",
		"/org/freedesktop/portal/desktop",
		"org.freedesktop.portal.ScreenCast.SelectSources",
		0,
		handleToken,
		nil,
		session,
		vardict{
			"handle_token": dbus.MakeVariant(handleToken),
			// TODO: need to check if cursor mode is available first
			"cursor_mode": dbus.MakeVariant(dbusCursorModeEmbedded),
			// TODO: session persistence
		},
	)
}

func (p *Pipewire) Init() error {
	session, err := p.createSession()
	if err != nil {
		return fmt.Errorf("createSession: %w", err)
	}
	fmt.Println(session)
	if err := p.selectSources(session); err != nil {
		return fmt.Errorf("selectSources: %w", err)
	}
	return nil
}

func (p *Pipewire) Close() error {
	return p.dbusConn.Close()
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
