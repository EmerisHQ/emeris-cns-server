package rest

import (
	models "github.com/allinbits/demeris-backend-models/cns"
	"net/http"

	"github.com/gin-gonic/gin"
)

const getChainsRoute = "/chains"

type getChainsResp struct {
	Chains []models.Chain `json:"chains"`
}

// @Summary Retrieve all chains
// @Description Get all chains added to the CNS
// @Router /chains [get]
// @Produce json
// @Success 200 {object} getChainResp
// @Failure 500
func (r *router) getChainsHandler(ctx *gin.Context) {
	data, err := r.s.d.Chains()

	if err != nil {
		e(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getChainsResp{
		Chains: data,
	})
}
func (r *router) getChains() (string, gin.HandlerFunc) {
	return getChainsRoute, r.getChainsHandler
}
