// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
	"github.com/edgexfoundry/device-virtual-go"
	"github.com/edgexfoundry/device-virtual-go/driver"
)

const (
	version     string = device_virtual.Version
	serviceName string = "device-virtual"
)

func main() {
	d := driver.NewVirtualDeviceDriver()
	startup.Bootstrap(serviceName, version, d)
}
