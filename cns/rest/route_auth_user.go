package rest

import (
	"errors"
	"net/http"

	"github.com/allinbits/emeris-cns-server/cns/auth"
	"github.com/gin-gonic/gin"
)

type userInfoResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *router) User(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		e(ctx, http.StatusBadRequest, errors.New("no user"))
		r.s.l.Error("no user")
		return
	}

	ctx.JSON(http.StatusOK, userInfoResponse{
		user.(auth.User).Name,
		user.(auth.User).Email,
	})

}
