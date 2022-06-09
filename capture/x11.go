//go:build linux || freebsd || openbsd || dragonfly

package capture

import (
	"image"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xfixes"
	"github.com/jezek/xgb/xproto"
)

// See https://www.ssec.wisc.edu/~billh/bp/xshm.c for a minimal implementation
// of this strategy.

// X11 uses the X library for screen capture.
type X11 struct {
	conn  *xgb.Conn
	setup *xproto.SetupInfo
	root  xproto.Window
}

func NewX11() (ret *X11, err error) {
	ret = &X11{}

	ret.conn, err = xgb.NewConn()
	if err != nil {
		return nil, err
	}
	ret.setup = xproto.Setup(ret.conn)
	ret.root = ret.setup.Roots[0].Root

	if err := xfixes.Init(ret.conn); err != nil {
		return nil, err
	}
	if err := xfixes.SelectCursorInputChecked(ret.conn, ret.root, xfixes.CursorNotifyMaskDisplayCursor).Check(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (x *X11) Close() {
	x.conn.Close()
}

func (x *X11) Start(framerate uint32, rect image.Rectangle, cb func(image.Image)) {
	go func() {
		_ = x.setup.DefaultScreen(x.conn)
	}()
}
