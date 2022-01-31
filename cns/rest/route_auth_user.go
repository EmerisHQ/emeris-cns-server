package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type userInfoResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *router) User(ctx *gin.Context) {

	name, _ := ctx.Get("name")
	email, _ := ctx.Get("email")

	ctx.JSON(http.StatusOK, userInfoResponse{
		name.(string),
		email.(string),
	})

}
