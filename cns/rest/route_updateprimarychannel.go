package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const updatePrimaryChannelRoute = "/update_primary_channel"

type updatePrimaryChannelRequest struct {
	Chain          string `json:"chain_name"`
	DestChain      string `json:"dest_chain"`
	PrimaryChannel string `json:"primary_channel"`
}

// @Summary Update the primary channel
// @Description Update the primary channel for a chain
// @Router /update_primary_channel [post]
// @Param chain body updatePrimaryChannelRequest true "Chain data to update"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 "if cannot parse payload"
// @Failure 500
func (r *router) updatePrimaryChannelHandler(ctx *gin.Context) {
	req := updatePrimaryChannelRequest{}

	if err := ctx.BindJSON(&req); err != nil {
		e(ctx, http.StatusBadRequest, err)
		r.s.l.Error("cannot bind json to updatePrimaryChannelRequest", err)
		return

	}

	if err := r.s.DB.UpdatePrimaryChannel(req.Chain, req.DestChain, req.PrimaryChannel); err != nil {
		e(ctx, http.StatusInternalServerError, err)
		r.s.l.Error("cannot update primary channel", err)
		return
	}

	return
}
func (r *router) updatePrimaryChannel() (string, gin.HandlerFunc) {
	return updatePrimaryChannelRoute, r.updatePrimaryChannelHandler
}
