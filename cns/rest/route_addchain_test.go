package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/allinbits/demeris-backend-models/cns"
	"github.com/allinbits/emeris-cns-server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"testing"
)

// Test adding a Chain with & without PublicEndpointInfo
// Mocks the K8s client
func TestAddChain(t *testing.T) {
	t.Parallel()

	// prepare the mock K8s client
	setupMockK8sClient()

	tests := []struct {
		name        string
		dataStruct  cns.Chain
		expectedHttpCode int
		success bool
	}{
		{
			"Add Chain - Invalid",
			cns.Chain{},
			400,
			true,
		},
		{
			"Add Chain - Without PublicEndpoint",
			chainWithoutPublicEndpoints,
			201,
			true,
		},
		{
			"Add Chain - With PublicEndpoints",
			chainWithPublicEndpoints,
			201,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// arrange
			addReq := addChainRequest {
				tt.dataStruct,
				true,
				&nodeConfig,
				&relayerConfig,
			}
			payload, _ := json.Marshal(addReq)

			// act
			resp, err := http.Post(fmt.Sprintf("http://%s%s", testingCtx.server.config.RESTAddress, addChainRoute), "application/json", strings.NewReader(string(payload)))

			// assert
			if !tt.success {
				require.Error(t, err, "Expecting a failed test case")
			} else {
				require.NoError(t, err)

				assert.Equal(t, tt.expectedHttpCode, resp.StatusCode)

				require.NoError(t, err)
				// TODO: ID is auto-calculated, testnet chain-ids are auto-populated.
				//   How to verify the chain from the DB?
			}
		})
	}
}

// Prepare the mock for the calls to expect
func setupMockK8sClient() {

	kubeClient := *testingCtx.server.k

	// Mocked 'List' does not "find" matching nodes (i.e. leave NodeSetList empty)
	kubeClient.(*mocks.Client).On("List",
		mock.Anything, 	// *context.emptyCtx
		mock.Anything, 	// *v1.NodeSetList
		mock.Anything,	// client.MatchingFields
		mock.Anything,	// client.InNamespace
		).Return(func(context.Context, client.ObjectList, ...client.ListOption) error { return nil })

	// Mocked 'Create' does nothing
	kubeClient.(*mocks.Client).On("Create",
		mock.Anything, 	// context.Context
		mock.Anything, 	// client.ObjectList
	).Return(func(context.Context, client.Object, ...client.CreateOption) error { return nil })
}
