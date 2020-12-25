package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/inahga/acolyte/pkg/drm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <path to gpu>\n", os.Args[0])
		os.Exit(2)
	}

	card, err := drm.Open(os.Args[1])
	if err != nil {
		panic(fmt.Errorf("open: %s", err))
	}
	defer card.Close()

	type connector struct {
		Connector  *drm.ModeConnector
		Properties []*drm.ModeProperty
	}
	dump := struct {
		Version      *drm.Version
		Resources    *drm.ModeResources
		CRTCs        []*drm.ModeCRTC
		Encoders     []*drm.ModeEncoder
		Connectors   []connector
		Blobs        []*drm.ModeBlob
		Planes       []*drm.ModePlane
		Framebuffers []*drm.ModeFramebuffer
	}{}

	ver, err := card.Version()
	if err != nil {
		panic(fmt.Errorf("version: %s", err))
	}
	dump.Version = ver

	for _, cap := range []uint64{drm.ClientCapAtomic, drm.ClientCapUniversalPlanes,
		drm.ClientCapWritebackConnectors} {
		if err := card.SetClientCap(cap, 1); err != nil {
			panic(fmt.Errorf("setcap: %w", err))
		}
	}

	res, err := card.ModeGetResources()
	if err != nil {
		panic(fmt.Errorf("resources: %s", err))
	}
	dump.Resources = res

	for _, conn := range res.ConnectorIDs {
		c, err := card.ModeGetConnector(conn)
		if err != nil {
			panic(fmt.Errorf("connector: %s", err))
		}

		var properties []*drm.ModeProperty
		for index, prop := range c.PropIDs {
			p, err := card.ModeGetProperty(prop)
			if err != nil {
				panic(fmt.Errorf("property: %s", err))
			}
			properties = append(properties, p)

			if p.Flags&drm.ModePropBlob != 0 && c.PropValues[index] != 0 {
				b, err := card.ModeGetBlob(uint32(c.PropValues[index]))
				if err != nil {
					panic(fmt.Errorf("blob: %s", err))
				}
				dump.Blobs = append(dump.Blobs, b)
			}
		}

		dump.Connectors = append(dump.Connectors, connector{
			Connector:  c,
			Properties: properties,
		})
	}

	for _, crtc := range res.CRTCIDs {
		c, err := card.ModeGetCRTC(crtc)
		if err != nil {
			panic(fmt.Errorf("crtc: %s", err))
		}
		dump.CRTCs = append(dump.CRTCs, c)
	}

	for _, encoder := range res.EncoderIDs {
		e, err := card.ModeGetEncoder(encoder)
		if err != nil {
			panic(fmt.Errorf("encoder: %s", err))
		}
		dump.Encoders = append(dump.Encoders, e)
	}

	planes, err := card.ModeGetPlaneResources()
	if err != nil {
		panic(fmt.Errorf("planes: %s", err))
	}
	for _, id := range *planes {
		plane, err := card.ModeGetPlane(id)
		if err != nil {
			panic(fmt.Errorf("plane: %s", err))
		}
		dump.Planes = append(dump.Planes, plane)
	}

	for _, fb := range res.FBIDs {
		f, err := card.ModeGetFramebuffer(fb)
		if err != nil {
			panic(fmt.Errorf("framebuffer: %s", err))
		}
		dump.Framebuffers = append(dump.Framebuffers, f)
	}

	b, err := json.MarshalIndent(dump, "", "    ")
	if err != nil {
		panic(fmt.Errorf("marshal: %s", err))
	}
	fmt.Println(string(b))
}
