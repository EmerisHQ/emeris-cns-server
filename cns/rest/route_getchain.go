package rest

import (
	"fmt"
	"net/http"
	"strings"

	models "github.com/allinbits/demeris-backend-models/cns"

	"github.com/gin-gonic/gin"
)

const (
	getChainRoute = "/chain/:chain"
	chainNotFoundErrorMsg = "no rows in result set"
)

type getChainResp struct {
	Chain models.Chain `json:"chain"`
}

// @Summary Retrieve a chain
// @Description get chain by name
// @Router /chain/{chain} [get]
// @Param chain path string true "Chain name to return"
// @Produce json
// @Success 200 {object} getChainResp
// @Failure 400 "if name is missing"
// @Failure 404 "if chain not found"
// @Failure 500 "on error"
func (r *router) getChainHandler(ctx *gin.Context) {

	chain, ok := ctx.Params.Get("chain")

	if !ok {
		e(ctx, http.StatusBadRequest, fmt.Errorf("chain not supplied"))
	}

	data, err := r.s.d.Chain(chain)

	if err != nil {
		if strings.Contains(err.Error(), chainNotFoundErrorMsg) {
			e(ctx, http.StatusNotFound, err)
		} else {
			e(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, getChainResp{
		Chain: data,
	})
}
func (r *router) getChain() (string, gin.HandlerFunc) {
	return getChainRoute, r.getChainHandler
}
