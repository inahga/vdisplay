package capture

import (
	"github.com/inahga/vdisplay/internal/drm"
)

// Writeback uses DRM writeback connectors for output.
//
// See https://gitlab.freedesktop.org/wayland/weston/-/merge_requests/458
// for example implementation.
//
// This capture backend is not possible without more widespread support for DRM
// leasing. The required mode and prop setting is gated by being the DRM master,
// which won't reasonably be us.
type Writeback struct {
	// card      *drm.Card
	// connector *drm.ModeConnector
}

func NewWriteback(card *drm.Card) (*Writeback, error) {
	panic("not supported")
	// ret := &Writeback{card: card}

	// resources, err := card.ModeGetResources()
	// if err != nil {
	// 	return nil, err
	// }
	// for _, id := range resources.ConnectorIDs {
	// 	candidate, err := card.ModeGetConnector(id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if candidate.Type == drm.ModeConnectorWriteback {
	// 		ret.connector = candidate
	// 		break
	// 	}
	// }
	// if ret.connector == nil {
	// 	return nil, fmt.Errorf("couldn't find writeback connector")
	// }

	// fmt.Println(card.ModeGetProperty(37))

	// return ret, nil
}
