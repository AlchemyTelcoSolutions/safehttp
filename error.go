package safehttp

import (
	"fmt"
	"net"
)

type SafeHTTPError struct {
	Err       error
	Text      string
	ErrCode   uint
	network   string
	address   string
	host      string
	ipAddress net.IP
}

var (
	BadNetworkTypeErrCode     uint = 1
	BadHostPortPairErrCode    uint = 2
	BadIPAddressErrCode       uint = 3
	NotPublicIPAddressErrCode uint = 4
)

func (err *SafeHTTPError) Unwrap() error {
	return err.Err
}

func (err *SafeHTTPError) Error() string {
	switch err.ErrCode {
	case BadNetworkTypeErrCode:
		return fmt.Sprintf("safehttp: %s is not a safe network type", err.network)
	case BadHostPortPairErrCode:
		return fmt.Sprintf("safehttp: %s is not a valid host/port pair", err.address)
	case BadIPAddressErrCode:
		return fmt.Sprintf("safehttp: %s is not a valid IP address", err.host)
	case NotPublicIPAddressErrCode:
		return fmt.Sprintf("safehttp: %s is not a public IP address", err.ipAddress)
	}
	return "safehttp: unknown error"
}

func newError(e error, code uint, network, address, host string, ipAddress net.IP) error {
	err := SafeHTTPError{
		Err:       e,
		ErrCode:   code,
		network:   network,
		address:   address,
		host:      host,
		ipAddress: ipAddress,
	}
	return &err
}
