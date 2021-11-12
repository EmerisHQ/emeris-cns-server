package rest

import (
	"encoding/json"
	"fmt"
	"github.com/allinbits/demeris-backend-models/cns"
	"github.com/google/go-cmp/cmp"
	"github.com/lib/pq"
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
			false,
		},
		{
			"Get Chain - Without PublicEndpoint",
			chainWithoutPublicEndpoints,
			chainWithoutPublicEndpoints.ChainName,
			200,
			true,
		},
		//{
		//	"Get Chain - With PublicEndpoints",
		//	chainWithPublicEndpoints,
		//	"chain2",
		//	200,
		//	true,
		//},
	}

	regExp := regexp.MustCompile(`/:\w+`)
	endpoint := regExp.ReplaceAllString(getChainRoute, "/%s")
	url := "http://%s" + endpoint

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// if we have a populated Chain store it
			if !cmp.Equal(tt.dataStruct, cns.Chain{}) {
				err := testingCtx.server.d.AddChain(tt.dataStruct)
				require.NoError(t, err)
			}

			// act
			resp, err := http.Get(fmt.Sprintf(url, testingCtx.server.config.RESTAddress, tt.chainName))
			defer func() { _ = resp.Body.Close() }()

			// assert
			if tt.success && err != nil {
				require.FailNow(t, "Unexpected test failure")
			} else if err != nil {
				require.NoError(t, err, "Unexpected error")
			} else {
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

var chainWithoutPublicEndpoints = cns.Chain {
	Enabled: true,
	ChainName: "chain1",
	Logo: "http://logo.com",
	DisplayName: "Chain 1",
	PrimaryChannel: map[string]string {"key": "value"},
	Denoms: []cns.Denom {
		cns.Denom {
			Name: "denom1",
			DisplayName: "Denom 1",
			Logo: "http://logo.com",
			Precision: 8,
			Verified: true,
			Stakable: true,
			Ticker: "DENOM1",
			PriceID: "price_id_1",
			FeeToken: true,
			GasPriceLevels: cns.GasPrice{
				Low: 0.2,
				Average: 0.3,
				High: 0.4,
			},
			FetchPrice: true,
			RelayerDenom: true,
			MinimumThreshRelayerBalance: nil,
		},
	},
	DemerisAddresses: []string {"12345"},
	GenesisHash: "hash",
	NodeInfo: cns.NodeInfo {
		Endpoint: "http://endpoint",
		ChainID: "chain_123",
		Bech32Config: cns.Bech32Config {
			MainPrefix: "prefix",
			PrefixAccount: "acc",
			PrefixValidator: "val",
			PrefixConsensus: "cons",
			PrefixPublic: "pub",
			PrefixOperator: "oper",
		},
	},
	ValidBlockThresh: cns.Threshold(30),
	DerivationPath: "m/44'/60'/0'/1",
	SupportedWallets: pq.StringArray([]string {"keplr"}),
	BlockExplorer: "http://explorer.com",
}

var chainWithPublicEndpoints = cns.Chain {

}