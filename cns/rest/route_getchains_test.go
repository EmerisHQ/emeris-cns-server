package rest

import (
	"encoding/json"
	"fmt"
	"github.com/allinbits/demeris-backend-models/cns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

// Test returning Chains with & without PublicEndpointInfo
func TestGetChains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		dataStructs  []cns.Chain
		expectedHttpCode int
		success bool
	}{
		{
			"Get Chains - Empty",
			[]cns.Chain{},	// ignored
			200,
			true,
		},
		{
			"Get Chains - One, without PublicEndpoint",
			[]cns.Chain{chainWithoutPublicEndpoints},
			200,
			true,
		},
		{
			"Get Chain - With PublicEndpoints",
			[]cns.Chain{chainWithoutPublicEndpoints, chainWithPublicEndpoints},
			200,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// if we have a populated Chain store it
			if len(tt.dataStructs) > 0 {
				for _, chain := range tt.dataStructs {
					err := testingCtx.server.d.AddChain(chain)
					require.NoError(t, err)
				}
			}

			// act
			resp, err := http.Get(fmt.Sprintf("http://%s%s", testingCtx.server.config.RESTAddress, getChainsRoute))
			defer func() { _ = resp.Body.Close() }()

			// assert
			if !tt.success {
				require.Error(t, err, "Expecting a failing test case")
				assert.Equal(t, tt.expectedHttpCode, resp.StatusCode)
			} else {
				require.NoError(t, err)

				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)

				respStruct := getChainsResp{}
				err = json.Unmarshal(body, &respStruct)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedHttpCode, resp.StatusCode)
				assert.Subset(t, respStruct.Chains, tt.dataStructs)
			}
		})
	}
}