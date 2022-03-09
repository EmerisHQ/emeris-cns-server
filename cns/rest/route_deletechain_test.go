package rest_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/emerishq/demeris-backend-models/cns"
	"github.com/google/go-cmp/cmp"

	"github.com/emerishq/emeris-cns-server/cns/rest"

	"github.com/stretchr/testify/require"
)

// Test deleting a Chain
func TestDeleteChain(t *testing.T) {

	tests := []struct {
		name             string
		dataStruct       cns.Chain
		chainName        string
		expectedHttpCode int
	}{
		{
			"Delete Chain - Unknown chain",
			cns.Chain{},
			"foo",
			404,
		},
		{
			"Delete Chain - Known chain",
			chainWithPublicEndpoints,
			"chain2",
			200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// arrange

			// if we have a populated Chain store, add it
			if !cmp.Equal(tt.dataStruct, cns.Chain{}) {
				err := testingCtx.server.DB.AddChain(tt.dataStruct)
				require.NoError(t, err)
			}

			payload, _ := json.Marshal(rest.DeleteChainRequest{Chain: tt.chainName})

			// act
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("http://%s%s", testingCtx.server.Config.RESTAddress, rest.DeleteChainRoute), strings.NewReader(string(payload)))
			req.Header.Add("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)

			// assert
			require.NoError(t, err)

			require.Equal(t, tt.expectedHttpCode, resp.StatusCode)

			require.NoError(t, err)
		})
	}
	truncateDB(t)
}
