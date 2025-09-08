package safehttp_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/AlchemyTelcoSolutions/safehttp"
	"github.com/stretchr/testify/assert"
)

func TestGetSafeHTTPCLient(t *testing.T) {
	type testCase struct {
		name                         string
		requestMethod                string
		requestURL                   string
		opts                         safehttp.Options
		isErrExpected                bool
		isSafeHTTPErrorExpected      bool
		expectedSafeHTTPErrorCode    uint
		expectedSafeHTTPErrorMessage string
	}
	cases := []testCase{
		{
			name:                         "simple localhost",
			opts:                         safehttp.Options{},
			requestMethod:                "GET",
			requestURL:                   "http://localhost:5000",
			isErrExpected:                true,
			isSafeHTTPErrorExpected:      true,
			expectedSafeHTTPErrorCode:    safehttp.NotPublicIPAddressErrCode,
			expectedSafeHTTPErrorMessage: "is not a public IP address",
		},
		{
			name:                         "localhost as ipv4",
			opts:                         safehttp.Options{},
			requestMethod:                "GET",
			requestURL:                   "http://127.0.0.1:5000",
			isErrExpected:                true,
			isSafeHTTPErrorExpected:      true,
			expectedSafeHTTPErrorCode:    safehttp.NotPublicIPAddressErrCode,
			expectedSafeHTTPErrorMessage: "safehttp: 127.0.0.1 is not a public IP address",
		},
		{
			name:                         "localhost as ipv6",
			opts:                         safehttp.Options{},
			requestMethod:                "GET",
			requestURL:                   "http://0:0:0:0:0:0:0:1:5000",
			isErrExpected:                true,
			isSafeHTTPErrorExpected:      true,
			expectedSafeHTTPErrorCode:    safehttp.NotPublicIPAddressErrCode,
			expectedSafeHTTPErrorMessage: "safehttp: ::1 is not a public IP address",
		},
		{
			name:                         "localhost as short ipv6",
			opts:                         safehttp.Options{},
			requestMethod:                "GET",
			requestURL:                   "http://::1:5000",
			isErrExpected:                true,
			isSafeHTTPErrorExpected:      true,
			expectedSafeHTTPErrorCode:    safehttp.NotPublicIPAddressErrCode,
			expectedSafeHTTPErrorMessage: "safehttp: ::1 is not a public IP address",
		},
		{
			name: "localhost as ipv6 but ipv6 not allowed",
			opts: safehttp.Options{
				AllowedNetworkTypes: []string{"tcp4"},
			},
			requestMethod:                "GET",
			requestURL:                   "http://0:0:0:0:0:0:0:1:5000",
			isErrExpected:                true,
			isSafeHTTPErrorExpected:      true,
			expectedSafeHTTPErrorCode:    safehttp.BadNetworkTypeErrCode,
			expectedSafeHTTPErrorMessage: "safehttp: tcp6 is not a safe network type",
		},
		{
			name: "google but ip is not allowed",
			opts: safehttp.Options{
				AllowedNetworkTypes: []string{"tcp4"},
				ForbiddenIPs:        []string{"142.251.33.206"},
			},
			requestMethod:                "GET",
			requestURL:                   "http://142.251.33.206",
			isErrExpected:                true,
			isSafeHTTPErrorExpected:      true,
			expectedSafeHTTPErrorCode:    safehttp.BadIPAddressErrCode,
			expectedSafeHTTPErrorMessage: "safehttp: 142.251.33.206 is not a valid IP address",
		},
		{
			name:                    "regular errors are preserved",
			opts:                    safehttp.Options{},
			requestMethod:           "GET",
			requestURL:              "http://test.some.bogus.url:5000",
			isErrExpected:           true,
			isSafeHTTPErrorExpected: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// create client
			client := safehttp.NewClient(c.opts)
			assert.NotNil(t, client)
			// set bad local endpoint
			req, err := http.NewRequest(c.requestMethod, c.requestURL, nil)
			assert.Nil(t, err)
			// do actual request, should get a safehttp error
			res, err := client.Do(req)
			if !c.isErrExpected {
				assert.Nil(t, err)
				return
			}

			assert.NotNil(t, err)

			var safeHTTPError *safehttp.SafeHTTPError

			if !c.isSafeHTTPErrorExpected {
				if errors.As(err, &safeHTTPError) {
					assert.Fail(t, "error should not be a safe http error")
				}
				return
			}
			assert.ErrorAs(t, err, &safeHTTPError)
			assert.Contains(t, safeHTTPError.Error(), c.expectedSafeHTTPErrorMessage)
			assert.Equal(t, c.expectedSafeHTTPErrorCode, safeHTTPError.ErrCode)
			assert.Nil(t, res)
		})
	}
}
