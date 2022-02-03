package rest

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func (r *router) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if r.s.a.Env == "test" {
			r.s.l.Infow("skipping auth in test")
			ctx.Set("email", "tester@tendermint.com")
			ctx.Set("name", "Test Ickle")
			ctx.Next()
		}

		if cookie, err := ctx.Request.Cookie("auth._token.google"); err == nil {

			reqTokenEncoded := cookie.Value

			reqToken, err := url.PathUnescape(reqTokenEncoded)

			if err != nil {
				e(ctx, http.StatusUnauthorized, err)
				r.s.l.Error("error unescaping token", err)
				return
			}

			splitToken := strings.Split(reqToken, " ")

			if len(splitToken) != 2 {
				e(ctx, http.StatusUnauthorized, errors.New("invalid token"))
				r.s.l.Error("invalid token")
				return
			}

			reqToken = splitToken[1]
			claims, err := r.s.a.ParseJWT(reqToken)

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
			e(ctx, http.StatusUnauthorized, err)
			r.s.l.Error("token not found")
			return
		}
	}
}
