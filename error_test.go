package safehttp

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnwrap(t *testing.T) {
	err := newError(errors.New("test"), 0, "", "", "", nil)
	var safeHTTPError *SafeHTTPError
	errors.As(err, &safeHTTPError)
	err = safeHTTPError.Unwrap()
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "test")
}
func TestError(t *testing.T) {
	err := SafeHTTPError{ErrCode: 1, network: "test", host: "test", ipAddress: net.IPv4(0, 0, 0, 0)}
	assert.Equal(t, "safehttp: test is not a safe network type", err.Error())
	err.ErrCode = 2
	assert.Equal(t, "safehttp:  is not a valid host/port pair", err.Error())
	err.ErrCode = 3
	assert.Equal(t, "safehttp: test is not a valid IP address", err.Error())
	err.ErrCode = 4
	assert.Equal(t, "safehttp: 0.0.0.0 is not a public IP address", err.Error())
	err.ErrCode = 5
	assert.Equal(t, "safehttp: unknown error", err.Error())

}
