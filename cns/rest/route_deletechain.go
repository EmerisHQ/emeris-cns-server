package rest

import (
	"errors"
	"net/http"

	"github.com/allinbits/demeris-backend-models/validation"

	"github.com/allinbits/emeris-cns-server/utils/k8s"

	"github.com/gin-gonic/gin"
)

const deleteChainRoute = "/delete"

type deleteChainRequest struct {
	Chain string `json:"chain" binding:"required"`
}

// @Summary Delete a chain's configuration
// @Description Delete a chain identified by name
// @Router /delete [delete]
// @Param chain body deleteChainRequest true "Chain name to delete"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 "if cannot parse payload"
// @Failure 500
func (r *router) deleteChainHandler(ctx *gin.Context) {
	chain := deleteChainRequest{}

	if err := ctx.ShouldBindJSON(&chain); err != nil {
		e(ctx, http.StatusBadRequest, validation.MissingFieldsErr(err, false))
		r.s.l.Error("cannot bind input data to Chain struct", err)
		return
	}

	k := k8s.Querier{Client: *r.s.KubeClient, Namespace: r.s.Config.KubernetesNamespace}

	if err := k.DeleteNode(chain.Chain); err != nil {
		// there isn't always a k8s nodeset for a given chain
		if !errors.Is(err, k8s.ErrNotFound) {
			e(ctx, http.StatusInternalServerError, err)
			r.s.l.Error("cannot delete chain", err)
			return
		}
	}

	if err := r.s.DB.DeleteChain(chain.Chain); err != nil {
		e(ctx, http.StatusInternalServerError, err)
		r.s.l.Error("cannot delete chain", err)
		return
	}

}
