package database

import "github.com/emerishq/emeris-utils/database"

const createDatabase = `
CREATE DATABASE IF NOT EXISTS cns;
`

const createTableChains = `
CREATE TABLE IF NOT EXISTS cns.chains (
	id serial unique primary key,
	enabled boolean default false,
	chain_name string not null,
	valid_block_thresh string not null,
	logo string not null,
	display_name string not null,
	primary_channel jsonb not null,
	denoms jsonb not null,
	demeris_addresses text[] not null,
	genesis_hash string not null,
	node_info jsonb not null,
	derivation_path string not null,
	unique(chain_name)
)
`

const deleteChain = `
DELETE FROM 
	cns.chains
WHERE
	chain_name=:chain_name; 
`

const insertChain = `
INSERT INTO cns.chains
	(
		chain_name,
		enabled,
		logo,
		display_name,
		valid_block_thresh,
		primary_channel,
		denoms,
		demeris_addresses,
		genesis_hash,
		node_info,
		derivation_path,
		supported_wallets,
		block_explorer,
		public_node_endpoints,
		cosmos_sdk_version
	)
VALUES
	(
		:chain_name,
		:enabled,
		:logo,
		:display_name,
		:valid_block_thresh,
		:primary_channel,
		:denoms,
		:demeris_addresses,
		:genesis_hash,
		:node_info,
		:derivation_path,
		:supported_wallets,
		:block_explorer,
		:public_node_endpoints,
		:cosmos_sdk_version
	)
ON CONFLICT
	(chain_name)
DO UPDATE SET 
		chain_name=EXCLUDED.chain_name, 
		enabled=EXCLUDED.enabled,
		valid_block_thresh=EXCLUDED.valid_block_thresh,
		logo=EXCLUDED.logo, 
		display_name=EXCLUDED.display_name, 
		primary_channel=EXCLUDED.primary_channel, 
		denoms=EXCLUDED.denoms, 
		demeris_addresses=EXCLUDED.demeris_addresses, 
		genesis_hash=EXCLUDED.genesis_hash,
		node_info=EXCLUDED.node_info,
		derivation_path=EXCLUDED.derivation_path,
		supported_wallets=EXCLUDED.supported_wallets,
		block_explorer=EXCLUDED.block_explorer,
		public_node_endpoints=EXCLUDED.public_node_endpoints,
		cosmos_sdk_version=EXCLUDED.cosmos_sdk_version;
`

const getAllChains = `
SELECT * FROM cns.chains
`

const getChain = `
SELECT * FROM cns.chains WHERE chain_name=:chain limit 1;
`

const channelsBetweenChains = `
SELECT
	c1.chain_name AS chain_a_chain_name,
	c1.channel_id AS chain_a_channel_id,
	c1.counter_channel_id AS chain_a_counter_channel_id,
	c1.chain_id AS chain_a_chain_id,
	c1.state AS chain_a_state,
	c2.chain_name AS chain_b_chain_name,
	c2.channel_id AS chain_b_channel_id,
	c2.counter_channel_id AS chain_b_counter_channel_id,
	c2.chain_id AS chain_b_chain_id,
	c2.state AS chain_b_state
FROM
	(
		SELECT
			tracelistener.channels.chain_name,
			tracelistener.channels.channel_id,
			tracelistener.channels.counter_channel_id,
			tracelistener.clients.chain_id,
			tracelistener.channels.state
		FROM
			tracelistener.channels
			LEFT JOIN tracelistener.connections ON
					tracelistener.channels.hops[1]
					= tracelistener.connections.connection_id
			LEFT JOIN tracelistener.clients ON
					tracelistener.clients.client_id
					= tracelistener.connections.client_id
		WHERE
			tracelistener.connections.chain_name
			= tracelistener.channels.chain_name
			AND tracelistener.clients.chain_name
				= tracelistener.channels.chain_name
	)
		AS c1,
	(
		SELECT
			tracelistener.channels.chain_name,
			tracelistener.channels.channel_id,
			tracelistener.channels.counter_channel_id,
			tracelistener.clients.chain_id,
			tracelistener.channels.state
		FROM
			tracelistener.channels
			LEFT JOIN tracelistener.connections ON
					tracelistener.channels.hops[1]
					= tracelistener.connections.connection_id
			LEFT JOIN tracelistener.clients ON
					tracelistener.clients.client_id
					= tracelistener.connections.client_id
		WHERE
			tracelistener.connections.chain_name
			= tracelistener.channels.chain_name
			AND tracelistener.clients.chain_name
				= tracelistener.channels.chain_name
	)
		AS c2
WHERE
	c1.channel_id = c2.counter_channel_id
	AND c1.counter_channel_id = c2.channel_id
	AND c1.chain_name != c2.chain_name
	AND c1.state = '3'
	AND c2.state = '3'
	AND c1.chain_name = :source
	AND c2.chain_name = :destination
	AND c2.chain_id = :chainID
`

const addColumnSupportedWallets = `
ALTER TABLE cns.chains ADD COLUMN IF NOT EXISTS supported_wallets text[] NOT NULL DEFAULT '{}';
`
const addColumnBlockExplorer = `
ALTER TABLE cns.chains ADD COLUMN IF NOT EXISTS block_explorer string;
`
const addColumnPublicNodeEndpoints = `
ALTER TABLE cns.chains ADD COLUMN IF NOT EXISTS public_node_endpoints jsonb DEFAULT '{}';
`
const addColumnCosmosSDKVersion = `
ALTER TABLE cns.chains ADD COLUMN IF NOT EXISTS cosmos_sdk_version string DEFAULT 'v0.42.10';
`
const updateTendermintRPCEndpoints = `
UPDATE cns.chains 
SET public_node_endpoints = json_set(public_node_endpoints, '{tendermint_rpc}'::string[],json_build_array(public_node_endpoints->'tendermint_rpc')) 
WHERE public_node_endpoints::jsonb ? 'tendermint_rpc' and json_typeof(public_node_endpoints->'tendermint_rpc') = 'string';
`

const updateCosmosAPIEndpoints = `
UPDATE cns.chains 
SET public_node_endpoints = json_set(public_node_endpoints, '{cosmos_api}'::string[],json_build_array(public_node_endpoints->'cosmos_api')) 
WHERE public_node_endpoints::jsonb ? 'cosmos_api' and json_typeof(public_node_endpoints->'cosmos_api') = 'string';
`

var migrationList = []string{
	createDatabase,
	createTableChains,
	addColumnSupportedWallets,
	addColumnBlockExplorer,
	addColumnPublicNodeEndpoints,
	addColumnCosmosSDKVersion,
	updateTendermintRPCEndpoints,
	updateCosmosAPIEndpoints,
}

func (i *Instance) runMigrations() {
	if err := database.RunMigrations(i.connString, migrationList); err != nil {
		panic(err)
	}
}
