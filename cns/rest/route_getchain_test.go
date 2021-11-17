package rest

import (
	"encoding/json"
	"fmt"
	"github.com/allinbits/demeris-backend-models/cns"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
)

// Test returning a Chain with / without PublicEndpointInfo
func TestGetChain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		dataStruct  cns.Chain
		chainName	string
		expectedHttpCode int
		success bool
	}{
		{
			"Get Chain - Unknown chain",
			cns.Chain{},	// ignored
			"foo",
			404,
			true,
		},
		{
			"Get Chain - Without PublicEndpoint",
			chainWithoutPublicEndpoints,
			chainWithoutPublicEndpoints.ChainName,
			200,
			true,
		},
		{
			"Get Chain - With PublicEndpoints",
			chainWithPublicEndpoints,
			chainWithPublicEndpoints.ChainName,
			200,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			regExp := regexp.MustCompile(`/:\w+`)
			endpoint := regExp.ReplaceAllString(getChainRoute, "/%s")
			url := "http://%s" + endpoint

			// if we have a populated Chain store it
			if !cmp.Equal(tt.dataStruct, cns.Chain{}) {
				err := testingCtx.server.d.AddChain(tt.dataStruct)
				require.NoError(t, err)
			}

			// act
			resp, err := http.Get(fmt.Sprintf(url, testingCtx.server.config.RESTAddress, tt.chainName))
			defer func() { _ = resp.Body.Close() }()

			// assert
			if !tt.success {
				require.Error(t, err, "Expecting a failed test case")
				assert.Equal(t, tt.expectedHttpCode, resp.StatusCode)
			} else {
				require.NoError(t, err)

				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)

				respStruct := getChainResp{}
				err = json.Unmarshal(body, &respStruct)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedHttpCode, resp.StatusCode)
				assert.Equal(t, tt.dataStruct, respStruct.Chain)
			}
		})
	}
}
