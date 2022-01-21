package database

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/allinbits/demeris-backend-models/cns"
	dbutils "github.com/allinbits/emeris-utils/database"
	"github.com/allinbits/emeris-utils/logging"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	defaultRPC = "http://localhost:26657"
	defaultAPI = "http://localhost:1317"
	chain1Name = "cosmos-hub"
	chain2Name = "akash"

	// insert two chains with empty public_node_endpoints
	insertChainsMigration = `INSERT INTO cns.chains 
			(
			id, enabled, chain_name, valid_block_thresh, logo, display_name, primary_channel, denoms, demeris_addresses, 
			genesis_hash, node_info, derivation_path, supported_wallets, block_explorer, public_node_endpoints
			) 
			VALUES 
			(
				1, true, '` + chain1Name + `', '10s', 'logo url', 'Cosmos Hub', 
				'{"cn1": "cn1", "cn2": "cn2"}',
				'[
					{"display_name":"STAKE","name":"stake","verified":true,"fee_token":true,"fetch_price":true,"fee_levels":{"low":1,"average":22,"high":42},"precision":6},
					{"display_name":"ATOM","name":"uatom","verified":true,"fetch_price":true,"fee_token":true,"fee_levels":{"low":1,"average":22,"high":42},"precision":6}
				]', 
				ARRAY['feeaddress'], 'genesis_hash', 
				'{"endpoint":"endpoint","chain_id":"chainid","bech32_config":{"main_prefix":"main_prefix","prefix_account":"prefix_account","prefix_validator":"prefix_validator",
				"prefix_consensus":"prefix_consensus","prefix_public":"prefix_public","prefix_operator":"prefix_operator"}}', 
				'm/44''/118''/0''/0/0', ARRAY['cosmostation'], 'mintscan','{}'
			),
			(
				2, true, '` + chain2Name + `', '10s', 'logo url', 'Akash Network', 
				'{"cn3": "cn3", "cn4": "cn4"}',
				'[
					{"display_name":"AKT","name":"uakt","verified":true,"fetch_price":true,"fee_token":true,"fee_levels":{"low":1,"average":22,"high":42},"precision":6}
				]', 
				ARRAY['feeaddress2'], 'genesis_hash_2', 
				'{"endpoint":"endpoint2","chain_id":"chainid2","bech32_config":{"main_prefix":"main_prefix","prefix_account":"prefix_account","prefix_validator":"prefix_validator",
				"prefix_consensus":"prefix_consensus","prefix_public":"prefix_public","prefix_operator":"prefix_operator"}}', 
				'm/44''/118''/0''/0/0',ARRAY['cosmostation'], 'mintscan','{}'
			);
			`
	testDBMigrations = []string{
		createDatabase,
		createTableChains,
		insertChainsMigration,
	}
)

// global variables for tests
var (
	dbInstance *Instance
	logger     *zap.SugaredLogger
)

func TestMain(m *testing.M) {
	// setup test DB
	var ts testserver.TestServer
	ts, dbInstance = setupDB()
	defer ts.Stop()

	// logger
	logger = logging.New(logging.LoggingConfig{
		LogPath: "",
		Debug:   true,
	})

	code := m.Run()
	os.Exit(code)
}

func setupDB() (testserver.TestServer, *Instance) {
	// start new cockroachDB test server
	ts, err := testserver.NewTestServer()
	checkNoError(err)

	err = ts.WaitForInit()
	checkNoError(err)

	// create new instance of db
	i, err := New(ts.PGURL().String())
	checkNoError(err)

	// create and insert data into db
	err = dbutils.RunMigrations(ts.PGURL().String(), testDBMigrations)
	checkNoError(err)

	return ts, i
}

func checkNoError(err error) {
	if err != nil {
		log.Fatalf("got error: %s", err)
	}
}

func TestPublicNodeEndpointsMigration(t *testing.T) {
	chains, err := dbInstance.Chains()
	require.NoError(t, err)
	require.Len(t, chains, 2)

	// updating one chain public_node_endpoints to include api, rpc endpoints in string format
	setValue := `{"tendermint_rpc":"` + defaultRPC + `", "cosmos_api":"` + defaultAPI + `"}`

	res, err := dbInstance.Instance.DB.Exec(fmt.Sprintf(`
	UPDATE cns.chains
	SET public_node_endpoints='%s'
	WHERE chain_name='%s'
	`, setValue, chain1Name))
	require.NoError(t, err)

	rows, _ := res.RowsAffected()
	require.Equal(t, rows, int64(1))

	// run default migrations
	dbInstance.runMigrations()

	tests := []struct {
		name      string
		chainName string
		expected  cns.PublicNodeEndpoints
	}{
		{
			"Public node endpoints updated to array",
			chain1Name,
			cns.PublicNodeEndpoints{
				TendermintRPC: []string{defaultRPC},
				CosmosAPI:     []string{defaultAPI},
			},
		},
		{
			"Empty public node endpoints",
			chain2Name,
			cns.PublicNodeEndpoints{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain, err := dbInstance.Chain(tt.chainName)

			require.NoError(t, err)
			require.Equal(t, tt.expected, chain.PublicNodeEndpoints)
		})
	}
}
