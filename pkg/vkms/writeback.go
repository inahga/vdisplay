package vkms

import (
	"fmt"

	"github.com/inahga/acolyte/pkg/drm"
)

func (c *Client) FindConnectorProperty(connector *drm.ModeConnector, name string) (*drm.ModeProperty, error) {
	for _, propID := range connector.PropIDs {
		prop, err := c.Card.ModeGetProperty(propID)
		if err != nil {
			return nil, err
		}
		if prop.Name == name {
			return prop, nil
		}
	}
	return nil, nil
}

func (c *Client) FindWritebackConnector(resources *drm.ModeResources) (*drm.ModeConnector, error) {
	for _, id := range resources.ConnectorIDs {
		connector, err := c.Card.ModeGetConnector(id)
		if err != nil {
			return nil, err
		}
		if connector.Type == drm.ModeConnectorWriteback {
			return connector, nil
		}
	}
	return nil, fmt.Errorf("no writeback connector found")
}

func FindActiveCRTC(card *drm.Card, resources *drm.ModeResources) (*drm.ModeCRTC, error) {
	// Unsure of correct logic, this seems fragile. Search all encoders for one
	// with valid CRTC ID. This will be the target CRTC ID.
	for _, id := range resources.EncoderIDs {
		encoder, err := card.ModeGetEncoder(id)
		if err != nil {
			return nil, err
		}
		if encoder.CRTCID > 0 {
			crtc, err := card.ModeGetCRTC(encoder.CRTCID)
			if err != nil {
				return nil, err
			}
			return crtc, nil
		}
	}
	return nil, fmt.Errorf("unable to determine the active CRTC")
}

func getCurrentCRTC(card *drm.Card, connector *drm.ModeConnector) (uint32, error) {
	for index, id := range connector.PropIDs {
		prop, err := card.ModeGetProperty(id)
		if err != nil {
			return 0, err
		}
		if prop.Name == "CRTC_ID" {
			return uint32(connector.PropValues[index]), nil
		}
	}
	return 0, fmt.Errorf("unable to determine current CRTC")
}

func addConnectorToCRTC(card *drm.Card, resources *drm.ModeResources, connector *drm.ModeConnector, crtc *drm.ModeCRTC) error {
	// Get any connectors already using this CRTC.
	var using []uint32
	for _, id := range resources.ConnectorIDs {
		conn, err := card.ModeGetConnector(id)
		if err != nil {
			return err
		}
		currentCRTC, err := getCurrentCRTC(card, conn)
		if err != nil {
			return err
		}
		if currentCRTC == crtc.ID {
			if connector.ID == conn.ID {
				// No action required.
				return nil
			}
			using = append(using, conn.ID)
		}
	}

	crtc.SetConnectors = append(using, connector.ID)
	// Hack to accomodate C weak typing. Kernel expects fb_id == -1 to choose
	// current framebuffer. fb_id is unsigned, so underflow manually.
	crtc.FBID = 0xFFFFFFFF
	return card.ModeSetCRTC(*crtc)
}

func prepareWriteback(card *drm.Card) error {
	// resources, err := card.ModeGetResources()
	// if err != nil {
	// 	return err
	// }

	// wbconn, err := findWritebackConnector(card, resources)
	// if err != nil {
	// 	return err
	// }

	// activeCRTC, err := FindActiveCRTC(card, resources)
	// if err != nil {
	// 	return err
	// }

	// return addConnectorToCRTC(card, resources, wbconn, activeCRTC)
	return nil
}
