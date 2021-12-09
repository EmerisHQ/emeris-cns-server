package mocks

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/allinbits/emeris-cns-server/cns/middleware"
	"github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

// Client is a mock type for the AuthClient type
type AuthClient struct {
	*mock.Mock
}

// Create provides a mock function with given fields
func (_m AuthClient) AuthUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := middleware.AuthHeader{}

		if err := ctx.ShouldBindHeader(&auth); err != nil {
			jsonErr := middleware.AuthError{
				Error: err.Error(),
			}

			_ = ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, jsonErr)
			return
		}

		jwtTokenHeader := strings.Split(auth.Token, " ")

		if len(jwtTokenHeader) != 2 || jwtTokenHeader[0] != "JWT" {

			err := fmt.Errorf("invalid auth token")

			jsonErr := middleware.AuthError{
				Error: err.Error(),
			}

			_ = ctx.Error(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, jsonErr)
			return
		}

		ctx.Next()
	}
}
