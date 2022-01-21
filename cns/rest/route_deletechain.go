package rest

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/allinbits/demeris-backend-models/validation"

	"github.com/gin-gonic/gin"
)

const DeleteChainRoute = "/delete"

type DeleteChainRequest struct {
	Chain string `json:"chain" binding:"required"`
}

// @Summary Delete a chain's configuration
// @Description Delete a chain identified by name
// @Router /delete [delete]
// @Param chain body DeleteChainRequest true "Chain name to delete"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 "if cannot parse payload"
// @Failure 404 "if unknown chain name"
// @Failure 500
func (r *router) deleteChainHandler(ctx *gin.Context) {
	chain := DeleteChainRequest{}

	if err := ctx.ShouldBindJSON(&chain); err != nil {
		e(ctx, http.StatusBadRequest, validation.MissingFieldsErr(err, false))
		r.s.l.Error("cannot bind input data to Chain struct", err)
		return
	}

	_, err := r.s.DB.Chain(chain.Chain)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			e(ctx, http.StatusNotFound, err)
		} else {
			e(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	if err := r.s.DB.DeleteChain(chain.Chain); err != nil {
		e(ctx, http.StatusInternalServerError, err)
		r.s.l.Error("cannot delete chain", err)
		return
	}

}
