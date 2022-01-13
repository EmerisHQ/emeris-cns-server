package chainwatch

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	kube "sigs.k8s.io/controller-runtime/pkg/client"

	models "github.com/allinbits/demeris-backend-models/cns"
	"github.com/allinbits/emeris-cns-server/cns/database"
	"github.com/allinbits/emeris-utils/k8s"
	v1 "github.com/allinbits/starport-operator/api/v1"
)

type Instance struct {
	l                *zap.SugaredLogger
	k                kube.Client
	defaultNamespace string
	c                *Connection
	db               *database.Instance
	relayerDebug     bool
}

func New(
	l *zap.SugaredLogger,
	k kube.Client,
	defaultNamespace string,
	c *Connection,
	db *database.Instance,
	relayerDebug bool,
) *Instance {
	return &Instance{
		l:                l,
		k:                k,
		defaultNamespace: defaultNamespace,
		c:                c,
		db:               db,
		relayerDebug:     relayerDebug,
	}

}

func (i *Instance) Run() {
	for range time.Tick(1 * time.Second) {
		chains, err := i.c.Chains()
		if err != nil {
			i.l.Errorw("cannot get chains from redis", "error", err)
			continue
		}

		if chains == nil {
			continue
		}

		i.l.Debugw("chains in cache", "list", chains)

		for _, chain := range chains {
			chainStatus, found, err := i.c.ChainStatus(chain.Name)
			if err != nil {
				i.l.Errorw("cannot query chain status from redis at beginning of chains loop", "chainName", chain.Name, "error", err)
				continue
			}

			if !found {
				chainStatus = starting
				if err := i.c.SetChainStatus(chain.Name, chainStatus); err != nil {
					i.l.Errorw("cannot set new chain status in redis", "chainName", chain.Name, "error", err, "newStatus", starting.String())
					continue
				}
			}

			q := k8s.Querier{
				Client:    i.k,
				Namespace: i.defaultNamespace,
			}

			n, err := q.ChainByName(chain.Name)
			if err != nil {
				i.l.Errorw("cannot get chains from k8s", "error", err)
				continue
			}

			i.l.Debugw("chain status", "name", chain.Name, "status", chainStatus.String())

			switch chainStatus {
			case starting:
				if n.Status.Phase != v1.PhaseRunning {
					i.l.Debugw("chain not in running phase", "name", n.Name, "phase", n.Status.Phase)
					if err := i.c.SetChainStatus(chain.Name, starting); err != nil {
						i.l.Errorw("cannot set chain status in redis", "chainName", chain.Name, "error", err, "newStatus", starting.String())
						continue
					}
					continue
				}

				// chain is now in running phase
				if err := i.c.SetChainStatus(chain.Name, running); err != nil {
					i.l.Errorw("cannot set chain status in redis", "chainName", chain.Name, "error", err, "newStatus", running.String())
					continue
				}

				i.l.Debugw("chain status update", "name", chain.Name, "new_status", running.String())
			case running:
				if err := i.c.SetChainStatus(chain.Name, relayerConnecting); err != nil {
					i.l.Errorw("cannot set chain status in redis", "chainName", chain.Name, "error", err, "newStatus", relayerConnecting.String())
					continue
				}
			case relayerConnecting:
				relayer, err := q.Relayer()
				if err != nil {
					i.l.Errorw("cannot get relayer", "error", err)
					continue
				}

				amt, err := i.db.ChainAmount()
				if err != nil {
					i.l.Errorw("cannot get amount of chains", "error", err)
					continue
				}

				chainStatuses := relayer.Status.ChainStatuses

				if amt != len(chainStatuses) {
					continue // corner case where the chain gets added, previous chains are already connected, but the operator still reports "Running" because the
					// reconcile cycle didn't get up yet.
				}

				phase := relayer.Status.Phase
				if phase != v1.RelayerPhaseRunning {
					if amt == 1 {
						if err := i.c.SetChainStatus(chain.Name, done); err != nil {
							i.l.Errorw("cannot set chain status in redis", "chainName", chain.Name, "error", err, "newStatus", done.String())
							continue
						}
					}
					continue
				}

				if err := i.relayerFinished(chain, relayer); err != nil {
					i.l.Debugw("error while running relayerfinished", "error", err)
					continue
				}

				if err := i.c.SetChainStatus(chain.Name, done); err != nil {
					i.l.Errorw("cannot set chain status in redis", "chainName", chain.Name, "error", err, "newStatus", done.String())
					continue
				}
			case done:
				if err := i.c.RemoveChain(chain); err != nil {
					i.l.Errorw("cannot remove chain from redis", "error", err)
				}

				if err := i.c.DeleteChainStatus(chain.Name); err != nil {
					i.l.Errorw("cannot delete chain status in redis", "chainName", chain.Name, "error", err)
					continue
				}
			}

		}
	}
}

func (i *Instance) relayerFinished(chain Chain, relayer v1.Relayer) error {
	if err := i.setPrimaryChannel(chain, relayer); err != nil {
		return err
	}

	return nil
}

func (i *Instance) setPrimaryChannel(_ Chain, relayer v1.Relayer) error {
	chainsMap := map[string]models.Chain{}

	chains, err := i.db.Chains()
	if err != nil {
		return err
	}

	for _, chain := range chains {
		i.l.Debugw("chain read", "chainID", chain.NodeInfo.ChainID, "name", chain.ChainName)
		chainsMap[chain.NodeInfo.ChainID] = chain
	}

	result := i.updatePrimaryChannelForChain(chainsMap, relayer)

	for _, chain := range result {
		if err := i.db.AddChain(chain); err != nil {
			return fmt.Errorf("error while updating chain %s, %w", chain.ChainName, err)
		}
	}

	return nil
}

func (i *Instance) updatePrimaryChannelForChain(chainsMap map[string]models.Chain, relayer v1.Relayer) map[string]models.Chain {

	paths := relayer.Status.Paths
	for chainID, chain := range chainsMap {
		i.l.Debugw("iterating chainsmap", "chainID", chainID)

		for _, path := range paths {
			i.l.Debugw("iterating path", "path", path)

			if _, found := path[chainID]; !found {
				i.l.Debugw("skipping path since it's not related to me", "chain", chain.ChainName)
				continue
			}

			for counterpartyChainID, value := range path {
				i.l.Debugw("beginning of path iteration", "counterpartyChainID", counterpartyChainID, "chainID", chainID)
				if counterpartyChainID == chainID {
					i.l.Debugw("found ourselves", "chainID", chainID)
					continue
				}

				counterparty, found := chainsMap[counterpartyChainID]
				i.l.Debugw("found counterparty in chainsMap", "counterparty name", counterparty.ChainName, "found", found)

				if !found {
					i.l.Panicw("found counterparty chain which isn't in chainsMap", "chainsMap", chainsMap, "counterparty", counterpartyChainID)
				}

				i.l.Debugw("updating chain", "chain to be update", chainsMap[chainID].ChainName, "counterparty", counterparty.ChainName, "value", value.ChannelID)
				if _, ok := chain.PrimaryChannel[counterparty.ChainName]; ok {
					// don't overwrite a primary channel that was set before
					continue
				}

				chain.PrimaryChannel[counterparty.ChainName] = path[chainID].ChannelID
			}
		}

		i.l.Debugw("new primary channel struct", "data", chain.PrimaryChannel)

		chainsMap[chainID] = chain
	}

	return chainsMap
}
