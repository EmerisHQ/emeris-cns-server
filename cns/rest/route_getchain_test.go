package rest_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"

	"github.com/allinbits/emeris-cns-server/cns/rest"

	"github.com/allinbits/demeris-backend-models/cns"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

// Test returning a Chain with / without PublicEndpointInfo
func TestGetChain(t *testing.T) {

	tests := []struct {
		name             string
		dataStruct       cns.Chain
		chainName        string
		expectedHttpCode int
		success          bool
	}{
		{
			"Get Chain - Unknown chain",
			cns.Chain{}, // ignored
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

	regExp := regexp.MustCompile(`/:\w+`)
	endpoint := regExp.ReplaceAllString(rest.GetChainRoute, "/%s")
	url := "http://%s" + endpoint

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// if we have a populated Chain store, add it
			if !cmp.Equal(tt.dataStruct, cns.Chain{}) {
				err := testingCtx.server.DB.AddChain(tt.dataStruct)
				require.NoError(t, err)
			}

			// act
			resp, err := http.Get(fmt.Sprintf(url, testingCtx.server.Config.RESTAddress, tt.chainName))
			defer func() { _ = resp.Body.Close() }()

			// assert
			if !tt.success {
				require.Error(t, err, "Expecting a failed test case")
				require.Equal(t, tt.expectedHttpCode, resp.StatusCode)
			} else {
				require.NoError(t, err)

				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)

				respStruct := rest.GetChainResp{}
				err = json.Unmarshal(body, &respStruct)
				require.NoError(t, err)

				require.Equal(t, tt.expectedHttpCode, resp.StatusCode)
				require.Equal(t, tt.dataStruct, respStruct.Chain)
			}
		})
		truncateDB(t)
	}
}
