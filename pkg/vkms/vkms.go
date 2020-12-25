// Package vkms allows interacting with the vkms driver in userspace.
package vkms

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/inahga/acolyte/pkg/drm"
)

type Client struct {
	Card *drm.Card
}

const (
	vkmsDriverName = "vkms"
)

var (
	ErrNotGPU  = errors.New("target is not a GPU")
	ErrNotVKMS = errors.New("target is not a VKMS GPU")
)

func Find(searchPath string) (*Client, error) {
	dir, err := ioutil.ReadDir(searchPath)
	if err != nil {
		return nil, fmt.Errorf("readdir: %w", err)
	}

	var ret *Client
	for _, ent := range dir {
		if ent.Mode()&os.ModeDevice != 0 {
			candidate, err := Open(path.Join(searchPath, ent.Name()))
			if err == nil {
				ret = candidate
				break
			} else if !errors.Is(err, ErrNotGPU) && !errors.Is(err, ErrNotVKMS) {
				return nil, err
			}
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("no vkms gpus detected, is the vkms kernel module loaded?")
	}
	return ret, nil
}

func Open(path string) (*Client, error) {
	var card *drm.Card

	card, err := drm.Open(path)
	if err != nil {
		return nil, err
	}

	version, err := card.Version()
	if err != nil {
		return nil, ErrNotGPU
	}
	if version.Name != vkmsDriverName {
		return nil, ErrNotVKMS
	}

	if err := card.SetClientCap(drm.ClientCapAtomic, 1); err != nil {
		return nil, fmt.Errorf("setcap atomic: %w", err)
	}
	if err := card.SetClientCap(drm.ClientCapWritebackConnectors, 1); err != nil {
		return nil, fmt.Errorf("setcap writeback: %w", err)
	}
	if err := prepareWriteback(card); err != nil {
		return nil, fmt.Errorf("writeback: %w", err)
	}

	return &Client{Card: card}, nil
}

func (c *Client) Close() error {
	return c.Card.Close()
}
