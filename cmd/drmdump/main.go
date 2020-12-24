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
		Version    *drm.Version
		Resources  *drm.ModeResources
		Connectors []connector
	}{}

	ver, err := card.Version()
	if err != nil {
		panic(fmt.Errorf("version: %s", err))
	}
	dump.Version = ver

	if err := card.SetClientCap(drm.ClientCapAtomic, 1); err != nil {
		panic(fmt.Errorf("setcap atomic: %s", err))
	}
	if err := card.SetClientCap(drm.ClientCapWritebackConnectors, 1); err != nil {
		panic(fmt.Errorf("setcap writeback: %s", err))
	}

	res, err := card.ModeResources()
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
		for _, prop := range c.PropIDs {
			p, err := card.ModeGetProperty(prop)
			if err != nil {
				panic(fmt.Errorf("property: %s", err))
			}
			properties = append(properties, p)
		}

		dump.Connectors = append(dump.Connectors, connector{
			Connector:  c,
			Properties: properties,
		})
	}

	b, err := json.MarshalIndent(dump, "", "    ")
	if err != nil {
		panic(fmt.Errorf("marshal: %s", err))
	}
	fmt.Println(string(b))
}
