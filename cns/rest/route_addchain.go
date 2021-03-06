package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v12 "k8s.io/api/core/v1"

	v1 "github.com/allinbits/starport-operator/api/v1"
	models "github.com/emerishq/demeris-backend-models/cns"
	"github.com/emerishq/demeris-backend-models/validation"
	"github.com/emerishq/emeris-cns-server/cns/chainwatch"
	"github.com/emerishq/emeris-utils/k8s/operator"
)

const AddChainRoute = "/add"

type AddChainRequest struct {
	models.Chain

	SkipChannelCreation  bool                           `json:"skip_channel_creation"`
	NodeConfig           *operator.NodeConfiguration    `json:"node_config"`
	RelayerConfiguration *operator.RelayerConfiguration `json:"relayer_configuration"`
}

// @Summary Add a new chain configuration
// @Description Add a new chain to the CNS DB
// @Router /add [post]
// @Param chain body AddChainRequest true "Chain data to add"
// @Accept json
// @Produce json
// @Success 201
// @Failure 400 "if cannot parse payload, or cannot validate: fees, denoms, relayer config"
// @Failure 500
func (r *router) addChainHandler(ctx *gin.Context) {
	newChain := AddChainRequest{}

	usr, _ := ctx.Get("user")

	r.s.l.Debug(usr)

	if err := ctx.ShouldBindJSON(&newChain); err != nil {
		e(ctx, http.StatusBadRequest, validation.MissingFieldsErr(err, false))
		r.s.l.Error("cannot bind input data to Chain struct", err)
		return
	}

	if err := validateFees(newChain.Chain); err != nil {
		e(ctx, http.StatusBadRequest, err)
		r.s.l.Error("fee validation failed", err)
		return
	}

	if err := validateDenoms(newChain.Chain); err != nil {
		e(ctx, http.StatusBadRequest, err)
		r.s.l.Error("fee validation failed", err)
		return
	}

	if newChain.NodeConfig != nil {
		newChain.NodeConfig.Namespace = r.s.Config.KubernetesNamespace

		newChain.NodeConfig.Name = newChain.ChainName

		// we trust that TestnetConfig holds the real chain ID
		if newChain.NodeConfig.TestnetConfig != nil &&
			*newChain.NodeConfig.TestnetConfig.ChainId != newChain.NodeInfo.ChainID {
			newChain.NodeInfo.ChainID = *newChain.NodeConfig.TestnetConfig.ChainId
		}

		newChain.NodeConfig.TracelistenerDebug = r.s.Config.Debug

		node, err := operator.NewNode(*newChain.NodeConfig)
		if err != nil {
			e(ctx, http.StatusBadRequest, err)
			r.s.l.Error("cannot add chain", err)
			return
		}

		switch newChain.NodeConfig.DisableMinFeeConfig {
		case true:
			node.Spec.Config.Nodes.TraceStoreContainer.ImagePullPolicy = v12.PullNever
		default:
			minGasPriceVal := newChain.RelayerToken().GasPriceLevels.Low / 2
			minGasPricesStr := fmt.Sprintf("%v%s", minGasPriceVal, newChain.RelayerToken().Name)

			cfgOverride := v1.ConfigOverride{
				App: []v1.TomlConfigField{
					{
						Key: "minimum-gas-prices",
						Value: v1.TomlConfigFieldValue{
							String: &minGasPricesStr,
						},
					},
				},
			}
			node.Spec.Config.Nodes.ConfigOverride = &cfgOverride
		}

		hasFaucet := false
		if node.Spec.Init != nil {
			hasFaucet = node.Spec.Init.Faucet != nil
		}

		if newChain.RelayerConfiguration == nil {
			newChain.RelayerConfiguration = &operator.DefaultRelayerConfiguration
		}

		if err := newChain.RelayerConfiguration.Validate(); err != nil {
			e(ctx, http.StatusBadRequest, err)
			r.s.l.Errorw("cannot validate relayer configuration", "error", err)
			return
		}

		if err := r.s.rc.AddChain(chainwatch.Chain{
			Name:                 newChain.ChainName,
			AddressPrefix:        newChain.NodeInfo.Bech32Config.MainPrefix,
			HasFaucet:            hasFaucet,
			SkipChannelCreation:  newChain.SkipChannelCreation,
			HDPath:               newChain.DerivationPath,
			RelayerConfiguration: *newChain.RelayerConfiguration,
		}); err != nil {
			e(ctx, http.StatusInternalServerError, err)
			r.s.l.Error("cannot add chain name to cache", err)
			return
		}
	}

	if err := r.s.DB.AddChain(newChain.Chain); err != nil {
		e(ctx, http.StatusInternalServerError, err)
		r.s.l.Error("cannot add chain", err)
		return
	}

	// return 201
	ctx.Writer.WriteHeader(http.StatusCreated)
}

func validateFees(c models.Chain) error {
	ft := c.FeeTokens()
	if len(ft) == 0 {
		return fmt.Errorf("no fee token specified")
	}

	for _, denom := range ft {
		if denom.GasPriceLevels.Empty() {
			return fmt.Errorf("fee levels for %s are not defined", denom.Name)
		}
	}

	return nil
}

func validateDenoms(c models.Chain) error {
	foundRelayerDenom := false
	for _, d := range c.Denoms {
		if d.RelayerDenom {
			if foundRelayerDenom {
				return fmt.Errorf("multiple relayer denoms detected")
			}

			if d.MinimumThreshRelayerBalance == nil {
				return fmt.Errorf("relayer denom detected but no relayer minimum threshold balance defined")
			}

			foundRelayerDenom = true
		}
	}

	return nil
}
