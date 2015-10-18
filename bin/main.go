package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/ucloud/docker-machine-ucloud"
)

func main() {
	plugin.RegisterDriver(ucloud.NewDriver("", ""))
}
