package main

import (
	"fmt"

	"github.com/inahga/acolyte/pkg/vkms"
)

func printObj(o interface{}) {
	fmt.Printf("%+v\n", o)
}

func main() {
	client, err := vkms.Find("/dev/dri")
	if err != nil {
		panic(fmt.Errorf("find: %s", err))
	}
	defer client.Close()

	if err := client.Card.SetMaster(); err != nil {
		panic(fmt.Errorf("setmaster: %s", err))
	}

	resources, err := client.Card.ModeGetResources()
	if err != nil {
		panic(err)
	}

	crtc, err := vkms.FindActiveCRTC(client.Card, resources)
	if err != nil {
		panic(err)
	}

	// I have no idea why a large size is necessary.
	// drivers/gpu/drm/drm_gem_framebuffer_helper.c:177 is the failure point.
	buf, err := client.Card.ModeCreateDumb(
		/*uint32(crtc.VDisplay)*/ 8192,
		/*uint32(crtc.HDisplay)*/ 8192,
		32,
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Card.ModeDestroyDumb(buf.Handle); err != nil {
			panic(err)
		}
	}()
	printObj(buf)

	fb, err := client.Card.ModeAddFramebuffer(
		uint32(crtc.VDisplay),
		uint32(crtc.HDisplay),
		buf.Pitch,
		32, 24, buf.Handle,
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Card.ModeRemoveFramebuffer(fb.ID); err != nil {
			panic(err)
		}
	}()
	printObj(fb)

	wbconn, err := client.FindWritebackConnector(resources)
	if err != nil {
		panic(err)
	}
	printObj(wbconn.ID)

	fbprop, err := client.FindConnectorProperty(wbconn, "WRITEBACK_FB_ID")
	if err != nil {
		panic(err)
	}
	printObj(fbprop.PropID)

	if err := client.Card.ModeConnectorSetProperty(wbconn.ID, fbprop.PropID, uint64(fb.ID)); err != nil {
		panic(err)
	}
}
