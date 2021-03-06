package database

import (
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"

	models "github.com/emerishq/demeris-backend-models/cns"
	dbutils "github.com/emerishq/emeris-utils/database"
)

type Instance struct {
	Instance   *dbutils.Instance
	connString string
}

func New(connString string) (*Instance, error) {
	i, err := dbutils.New(connString)
	if err != nil {
		return nil, err
	}

	ii := &Instance{
		Instance:   i,
		connString: connString,
	}

	ii.runMigrations()
	return ii, nil
}

func (i *Instance) AddChain(chain models.Chain) error {
	n, err := i.Instance.DB.PrepareNamed(insertChain)
	if err != nil {
		return err
	}

	res, err := n.Exec(chain)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("database delete statement had no effect")
	}

	err = n.Close()
	if err != nil {
		return nil
	}

	return nil
}

func (i *Instance) DeleteChain(chain string) error {
	n, err := i.Instance.DB.PrepareNamed(deleteChain)
	if err != nil {
		return err
	}

	res, err := n.Exec(map[string]interface{}{
		"chain_name": chain,
	})

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("database delete statement had no effect")
	}

	err = n.Close()
	if err != nil {
		return nil
	}

	return nil
}

func (i *Instance) Chain(chain string) (models.Chain, error) {
	var c models.Chain

	n, err := i.Instance.DB.PrepareNamed(getChain)
	if err != nil {
		return models.Chain{}, err
	}
	err = n.Get(&c, map[string]interface{}{
		"chain": chain,
	})

	return c, err
}

func (i *Instance) Chains() ([]models.Chain, error) {
	var c []models.Chain

	return c, i.Instance.Exec(getAllChains, nil, &c)
}

func (i *Instance) UpdatePrimaryChannel(sourceChain, destChain, channel string) error {

	res, err := i.Instance.DB.Exec(fmt.Sprintf(`
	UPDATE cns.chains
	SET primary_channel = primary_channel || jsonb_build_object('%s' , '%s')
	WHERE chain_name='%s'
	`, destChain, channel, sourceChain))

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("update failed")
	}

	return nil
}

func (i *Instance) GetDenoms(chain string) (models.DenomList, error) {

	var l models.DenomList

	return l, i.Instance.Exec("select json_array_elements(denoms) from cns.chains where chain_name=:chain;", map[string]interface{}{
		"chain": chain,
	}, &l)
}

func (i *Instance) UpdateDenoms(chain string, denoms models.DenomList) error {

	b, err := json.Marshal(denoms)

	if err != nil {
		return err
	}

	res, err := i.Instance.DB.Exec(fmt.Sprintf(`
	UPDATE cns.chains
	SET denoms = '%s'::jsonb
	WHERE chain_name='%s'
	`, string(b), chain))

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("update failed")
	}

	return nil
}

type channelsBetweenChain struct {
	ChainAName             string `db:"chain_a_chain_name"`
	ChainAChannelID        string `db:"chain_a_channel_id"`
	ChainACounterChannelID string `db:"chain_a_counter_channel_id"`
	ChainAChainID          string `db:"chain_a_chain_id"`
	ChainAState            int    `db:"chain_a_state"`
	ChainBName             string `db:"chain_b_chain_name"`
	ChainBChannelID        string `db:"chain_b_channel_id"`
	ChainBCounterChannelID string `db:"chain_b_counter_channel_id"`
	ChainBChainID          string `db:"chain_b_chain_id"`
	ChainBState            int    `db:"chain_b_state"`
}

func (i *Instance) ChannelsBetweenChains(source, destination, chainID string) (map[string]string, error) {

	var c []channelsBetweenChain

	n, err := i.Instance.DB.PrepareNamed(channelsBetweenChains)
	if err != nil {
		return map[string]string{}, err
	}

	if err := n.Select(&c, map[string]interface{}{
		"source":      source,
		"destination": destination,
		"chainID":     chainID,
	}); err != nil {
		return map[string]string{}, err
	}

	ret := map[string]string{}

	for _, cc := range c {
		// channel ID destination => channel ID on source
		ret[cc.ChainAChannelID] = cc.ChainBChannelID
	}

	err = n.Close()
	if err != nil {
		return nil, nil
	}

	return ret, nil
}

func (i *Instance) ChainAmount() (int, error) {
	var ret int
	return ret, i.Instance.DB.Get(&ret, "select count(id) from cns.chains")
}
