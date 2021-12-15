package rest_test

import (
	"github.com/allinbits/demeris-backend-models/cns"
	"github.com/allinbits/emeris-cns-server/utils/k8s/operator"
	v1 "github.com/allinbits/starport-operator/api/v1"
	"github.com/lib/pq"
)

var relayerBalance = int64(30000)

var chainWithoutPublicEndpoints = cns.Chain{
	Enabled:        true,
	ChainName:      "chain1",
	Logo:           "http://logo.com",
	DisplayName:    "Chain 1",
	PrimaryChannel: map[string]string{"key": "value"},
	Denoms: []cns.Denom{
		{
			Name:        "denom1",
			DisplayName: "Denom 1",
			Logo:        "http://logo.com",
			Precision:   8,
			Verified:    true,
			Stakable:    true,
			Ticker:      "DENOM1",
			PriceID:     "price_id_1",
			FeeToken:    true,
			GasPriceLevels: cns.GasPrice{
				Low:     0.2,
				Average: 0.3,
				High:    0.4,
			},
			FetchPrice:                  true,
			RelayerDenom:                true,
			MinimumThreshRelayerBalance: &relayerBalance,
		},
	},
	DemerisAddresses: []string{"12345"},
	GenesisHash:      "hash",
	NodeInfo: cns.NodeInfo{
		Endpoint: "http://endpoint",
		ChainID:  "chain_123",
		Bech32Config: cns.Bech32Config{
			MainPrefix:      "prefix",
			PrefixAccount:   "acc",
			PrefixValidator: "val",
			PrefixConsensus: "cons",
			PrefixPublic:    "pub",
			PrefixOperator:  "oper",
		},
	},
	ValidBlockThresh: cns.Threshold(30),
	DerivationPath:   "m/44'/60'/0'/1",
	SupportedWallets: pq.StringArray([]string{"keplr"}),
	BlockExplorer:    "http://explorer.com",
	CosmosSDKVersion: "v0.42.10",
}

var chainWithPublicEndpoints = cns.Chain{
	Enabled:        true,
	ChainName:      "chain2",
	Logo:           "http://logo.com",
	DisplayName:    "Chain 2",
	PrimaryChannel: map[string]string{"key": "value"},
	Denoms: []cns.Denom{
		{
			Name:        "denom2",
			DisplayName: "Denom 2",
			Logo:        "http://logo.com",
			Precision:   8,
			Verified:    true,
			Stakable:    true,
			Ticker:      "DENOM2",
			PriceID:     "price_id_1",
			FeeToken:    true,
			GasPriceLevels: cns.GasPrice{
				Low:     0.2,
				Average: 0.3,
				High:    0.4,
			},
			FetchPrice:                  true,
			RelayerDenom:                true,
			MinimumThreshRelayerBalance: &relayerBalance,
		},
	},
	DemerisAddresses: []string{"12345"},
	GenesisHash:      "hash",
	NodeInfo: cns.NodeInfo{
		Endpoint: "http://endpoint",
		ChainID:  "chain_123",
		Bech32Config: cns.Bech32Config{
			MainPrefix:      "prefix",
			PrefixAccount:   "acc",
			PrefixValidator: "val",
			PrefixConsensus: "cons",
			PrefixPublic:    "pub",
			PrefixOperator:  "oper",
		},
	},
	ValidBlockThresh: cns.Threshold(30),
	DerivationPath:   "m/44'/60'/0'/1",
	SupportedWallets: pq.StringArray([]string{"keplr"}),
	BlockExplorer:    "http://explorer.com",
	PublicNodeEndpoints: cns.PublicNodeEndpoints{
		TendermintRPC: "https://www.host.com:1234",
		CosmosAPI:     "https://host.foo.bar:2345",
	},
	CosmosSDKVersion: "v0.44.3",
}

var nodeRPC = v1.NodeRPC{
	Address: "host.com",
}

var joinConfig = v1.JoinConfig{
	Genesis: v1.GenesisDownload{
		FromNodeRPC: &nodeRPC,
	},
	Seeds: []v1.Peer{
		{
			Id:      "test",
			Address: "127.0.0.1",
		},
	},
	PersistentPeers: []v1.Peer{
		{
			Id:      "test",
			Address: "127.0.0.1",
		},
	},
}

var amount = "1234"
var testnetName = "testnet"

var validatorInitConfig = v1.ValidatorInitConfig{
	ChainId:     &testnetName,
	StakeAmount: &amount,
	Assets:      []string{"asset1", "asset2"},
}

var nodeConfig = operator.NodeConfiguration{
	Name:                "test",
	CLIName:             "gaiad",
	JoinConfig:          &joinConfig,
	TestnetConfig:       &validatorInitConfig,
	DockerImage:         "docker/image",
	DockerImageVersion:  "v1.0.0",
	TracelistenerImage:  "tracelistener/image",
	DisableMinFeeConfig: true,
	TracelistenerDebug:  false,
}

var relayerConfig = operator.RelayerConfiguration{
	MaxMsgNum:      1234,
	ClockDrift:     "10s",
	MaxGas:         345,
	TrustingPeriod: "10min",
}
