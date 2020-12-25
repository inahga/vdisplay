package main

import (
	"fmt"

	"github.com/inahga/acolyte/pkg/vkms"
)

func main() {
	client, err := vkms.Find("/dev/dri")
	if err != nil {
		panic(fmt.Errorf("find: %s", err))
	}
	defer client.Close()
}
