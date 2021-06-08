package safehttp

import (
	"net"
	"net/http"
	"syscall"
	"time"
)

type Options struct {
	AllowedNetworkTypes []string
	ForbiddenIPs        []string
}

func NewClient(opts Options) *http.Client {
	opts = getDefaultOpts(opts)

	safeDialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
		Control:   getSafeSocketControlFunc(opts),
	}

	safeTransport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           safeDialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	safeClient := &http.Client{
		Transport: safeTransport,
	}
	return safeClient
}

func getDefaultOpts(opts Options) Options {
	if opts.AllowedNetworkTypes == nil || len(opts.AllowedNetworkTypes) == 0 {
		opts.AllowedNetworkTypes = []string{"tcp4", "tcp6"}
	}
	return opts
}

func isAllowedNetwork(network string, allowedNetworks []string) bool {
	foundAllowed := false
	for _, allowedNetwork := range allowedNetworks {
		if network == allowedNetwork {
			foundAllowed = true
		}
	}
	return foundAllowed
}

func isForbiddenIPAddress(address string, forbiddenAddresses []string) bool {
	for _, forbiddenAddress := range forbiddenAddresses {
		if address == forbiddenAddress {
			return true
		}
	}
	return false
}

func getSafeSocketControlFunc(opts Options) func(network, address string, c syscall.RawConn) error {
	return func(network string, address string, conn syscall.RawConn) error {

		if !isAllowedNetwork(network, opts.AllowedNetworkTypes) {
			return newError(nil, BadNetworkTypeErrCode, network, "", "", nil)
		}

		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return newError(err, BadHostPortPairErrCode, network, address, host, nil)
		}

		ipaddress := net.ParseIP(host)
		if ipaddress == nil {
			return newError(err, BadIPAddressErrCode, network, address, host, ipaddress)
		}

		if isForbiddenIPAddress(ipaddress.String(), opts.ForbiddenIPs) {
			return newError(err, BadIPAddressErrCode, network, address, host, ipaddress)
		}

		if !isPublicIPAddress(ipaddress) {
			return newError(err, NotPublicIPAddressErrCode, network, address, host, ipaddress)
		}

		return nil
	}
}
