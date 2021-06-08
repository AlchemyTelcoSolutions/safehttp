package safehttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSafeSocketControlFuncBadHost(t *testing.T) {
	opts := getDefaultOpts(Options{})
	fn := getSafeSocketControlFunc(opts)
	err := fn("tcp4", "bad:address:", nil)
	assert.NotNil(t, err)
	var safeHTTPErr *SafeHTTPError
	assert.ErrorAs(t, err, &safeHTTPErr)
	assert.Equal(t, BadHostPortPairErrCode, safeHTTPErr.ErrCode)
}

func TestGetSafeSocketControlFuncBadIPAddress(t *testing.T) {
	opts := getDefaultOpts(Options{})
	fn := getSafeSocketControlFunc(opts)
	err := fn("tcp4", "12345789.12345789.123456789879:80", nil)
	assert.NotNil(t, err)
	var safeHTTPErr *SafeHTTPError
	assert.ErrorAs(t, err, &safeHTTPErr)
	assert.Equal(t, BadIPAddressErrCode, safeHTTPErr.ErrCode)
}

func TestGetSafeSocketControlFuncNoErr(t *testing.T) {
	opts := getDefaultOpts(Options{})
	fn := getSafeSocketControlFunc(opts)
	err := fn("tcp4", "142.251.33.206:80", nil)
	assert.Nil(t, err)
}
