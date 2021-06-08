package safehttp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIPv4Reserved(t *testing.T) {
	assert.False(t, isIPv4Reserved(net.IPv4(142, 251, 33, 206)))
	assert.True(t, isIPv4Reserved(net.IPv4(127, 0, 0, 1)))

}
