package rest

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	models "github.com/emerishq/demeris-backend-models/cns"

	"github.com/gin-gonic/gin"
)

const (
	GetChainRoute = "/chain/:chain"
)

type GetChainResp struct {
	Chain models.Chain `json:"chain"`
}

// @Summary Retrieve a chain
// @Description get chain by name
// @Router /chain/{chain} [get]
// @Param chain path string true "Chain name to return"
// @Produce json
// @Success 200 {object} GetChainResp
// @Failure 400 "if name is missing"
// @Failure 404 "if chain not found"
// @Failure 500 "on error"
func (r *router) getChainHandler(ctx *gin.Context) {

	chain, ok := ctx.Params.Get("chain")

	if !ok {
		e(ctx, http.StatusBadRequest, fmt.Errorf("chain not supplied"))
	}

	data, err := r.s.DB.Chain(chain)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			e(ctx, http.StatusNotFound, err)
		} else {
			e(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, GetChainResp{
		Chain: data,
	})
}
func (r *router) getChain() (string, gin.HandlerFunc) {
	return GetChainRoute, r.getChainHandler
}
