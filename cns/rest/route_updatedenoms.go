package rest

import (
	"net/http"

	models "github.com/allinbits/demeris-backend-models/cns"
	"github.com/gin-gonic/gin"
)

const updateDenomsRoute = "/denoms"

type updateDenomsRequest struct {
	Chain  string           `json:"chain_name"`
	Denoms models.DenomList `json:"denoms"`
}

// @Summary Update the denominations for a chain
// @Description Update the primary channel handler for a chain
// @Router /denoms [post]
// @Param chain body updateDenomsRequest true "Chain data to update"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 "if cannot parse payload"
// @Failure 500
func (r *router) updateDenomsHandler(ctx *gin.Context) {
	req := updateDenomsRequest{}

	if err := ctx.BindJSON(&req); err != nil {
		e(ctx, http.StatusBadRequest, err)
		r.s.l.Error("cannot bind json to updateDenomsRequest", err)
		return

	}

	if err := r.s.DB.UpdateDenoms(req.Chain, req.Denoms); err != nil {
		e(ctx, http.StatusInternalServerError, err)
		r.s.l.Error("cannot update denoms", err)
		return
	}

	return
}
func (r *router) updateDenoms() (string, gin.HandlerFunc) {
	return updateDenomsRoute, r.updateDenomsHandler
}
