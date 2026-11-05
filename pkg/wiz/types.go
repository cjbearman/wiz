package wiz

import (
	"time"
)

// jsonResponse represents a generic JSON response from the device
type jsonResponse interface {
	unmarshal(from []byte) error
	method() string
	id() int
	error() *deviceErrorReport
}

const (
	DefaultOperationTimeout = 1 * time.Second
	DefaultRetries          = 3
	pingTimeout             = 1 * time.Second
	pingCount               = 3
	pingInterval            = 100 * time.Millisecond
)

const (
	DefaultDiscoveryRunFor            = 2000 * time.Millisecond
	DefaultDiscoveryBroadcastInterval = 500 * time.Millisecond
	wizPort                           = 38899
)

type wizMethod string

const (
	getPilot        wizMethod = "getPilot"
	getSystemConfig wizMethod = "getSystemConfig"
	getUserConfig   wizMethod = "getUserConfig"
	getModelConfig  wizMethod = "getModelConfig"
	setPilot        wizMethod = "setPilot"
	setState        wizMethod = "setState"
	pulse           wizMethod = "pulse"
)
