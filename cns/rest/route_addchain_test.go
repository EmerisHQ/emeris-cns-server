package rest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/allinbits/emeris-cns-server/cns/rest"

	"github.com/allinbits/demeris-backend-models/cns"
	"github.com/allinbits/emeris-cns-server/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Test adding a Chain with & without PublicEndpointInfo
// Mocks the K8s client
func TestAddChain(t *testing.T) {

	setupMockK8sClient()

	tests := []struct {
		name             string
		dataStruct       cns.Chain
		expectedHttpCode int
		success          bool
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
			t.Parallel()

			// arrange
			addReq := rest.AddChainRequest{
				Chain:                tt.dataStruct,
				SkipChannelCreation:  true,
				NodeConfig:           &nodeConfig,
				RelayerConfiguration: &relayerConfig,
			}
			payload, _ := json.Marshal(addReq)

			// act
			resp, err := http.Post(fmt.Sprintf("http://%s%s", testingCtx.server.Config.RESTAddress, rest.AddChainRoute), "application/json", strings.NewReader(string(payload)))

			// assert
			if !tt.success {
				require.Error(t, err, "Expecting a failed test case")
			} else {
				require.NoError(t, err)

				require.Equal(t, tt.expectedHttpCode, resp.StatusCode)

				require.NoError(t, err)
				// TODO: ID is auto-calculated, testnet chain-ids are auto-populated.
				//   How to verify the chain from the DB?
			}
		})
	}
	truncateDB(t)
}

// Prepare the mock for the calls to expect
func setupMockK8sClient() {

	kubeClient := *testingCtx.server.KubeClient

	// Mocked 'List' does not "find" matching nodes (i.e. leave NodeSetList empty)
	kubeClient.(*mocks.Client).On("List",
		mock.Anything, // *context.emptyCtx
		mock.Anything, // *v1.NodeSetList
		mock.Anything, // client.MatchingFields
		mock.Anything, // client.InNamespace
	).Return(func(context.Context, client.ObjectList, ...client.ListOption) error { return nil })

	// Mocked 'Create' does nothing
	kubeClient.(*mocks.Client).On("Create",
		mock.Anything, // context.Context
		mock.Anything, // client.ObjectList
	).Return(func(context.Context, client.Object, ...client.CreateOption) error { return nil })
}
