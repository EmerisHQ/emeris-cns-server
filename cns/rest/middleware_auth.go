package rest

import (
	"net/http"
	"strings"

	"errors"

	"github.com/allinbits/emeris-cns-server/cns/auth"
	"github.com/gin-gonic/gin"
)

func (r *router) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if reqToken := ctx.Request.Header.Get("Authorization"); reqToken != "" {
			splitToken := strings.Split(reqToken, " ")
			reqToken = splitToken[1]
			claims, err := auth.ParseJWT(reqToken)

			if err != nil {
				e(ctx, http.StatusUnauthorized, err)
				r.s.l.Error("failed to verify token", err)
				return
			}

			email := claims["email"].(string)
			name := claims["name"].(string)
			ctx.Set("email", email)
			ctx.Set("name", name)

			r.s.l.Infow("incoming request from %s (%s)\n", name, email)
			ctx.Next()
		} else {
			e(ctx, http.StatusUnauthorized, errors.New("http://youtu.be/otCpCn0l4Wo?t=15"))
			r.s.l.Error("failed to verify token")
			return
		}
	}
}
