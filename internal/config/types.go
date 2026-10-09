// -*- mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2026 IOTech Ltd
// Copyright (c) 2021 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v4/config"
)

// WritableInfo is a struct which contains configuration settings that can be changed in the Registry .
type WritableInfo struct {
	// Level is the logging level of writing log message
	LogLevel        string
	InsecureSecrets config.InsecureSecrets
	Reading         Reading
	Telemetry       config.TelemetryInfo
}

// Reading is a struct which contains reading configuration settings.
type Reading struct {
	// ReadingUnits specifies whether or not to indicate the units of measure for the value in the reading
	ReadingUnits bool
}
