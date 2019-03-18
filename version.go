// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package device_virtual

import "io/ioutil"

// Global version for device-sdk-go
func GetVersion() string {
	if b, err := ioutil.ReadFile("../VERSION"); err != nil {
		return ""
	} else {
		return string(b)
	}
}
