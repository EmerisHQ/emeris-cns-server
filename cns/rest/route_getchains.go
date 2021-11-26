package rest

import (
	"net/http"

	models "github.com/allinbits/demeris-backend-models/cns"

	"github.com/gin-gonic/gin"
)

const GetChainsRoute = "/chains"

type GetChainsResp struct {
	Chains []models.Chain `json:"chains"`
}

// @Summary Retrieve all chains
// @Description Get all chains added to the CNS
// @Router /chains [get]
// @Produce json
// @Success 200 {object} GetChainResp
// @Failure 500
func (r *router) getChainsHandler(ctx *gin.Context) {
	data, err := r.s.DB.Chains()

	if err != nil {
		e(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, GetChainsResp{
		Chains: data,
	})
}
func (r *router) getChains() (string, gin.HandlerFunc) {
	return GetChainsRoute, r.getChainsHandler
}
