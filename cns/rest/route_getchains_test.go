package rest_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/emerishq/emeris-cns-server/cns/rest"

	"github.com/emerishq/demeris-backend-models/cns"
	"github.com/stretchr/testify/require"
)

// Test returning Chains with & without PublicEndpointInfo
func TestGetChains(t *testing.T) {

	tests := []struct {
		name             string
		dataStructs      []cns.Chain
		expectedHttpCode int
		success          bool
	}{
		{
			"Get Chains - Empty",
			[]cns.Chain{}, // ignored
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
					err := testingCtx.server.DB.AddChain(chain)
					require.NoError(t, err)
				}
			}

			// act
			resp, err := http.Get(fmt.Sprintf("http://%s%s", testingCtx.server.Config.RESTAddress, rest.GetChainsRoute))
			defer func() { _ = resp.Body.Close() }()

			// assert
			if !tt.success {
				require.Error(t, err, "Expecting a failing test case")
				require.Equal(t, tt.expectedHttpCode, resp.StatusCode)
			} else {
				require.NoError(t, err)

				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)

				respStruct := rest.GetChainsResp{}
				err = json.Unmarshal(body, &respStruct)
				require.NoError(t, err)

				require.Equal(t, tt.expectedHttpCode, resp.StatusCode)
				require.Subset(t, respStruct.Chains, tt.dataStructs)
			}
		})
		truncateDB(t)
	}
}
