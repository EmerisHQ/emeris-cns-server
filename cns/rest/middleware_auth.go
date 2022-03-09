package rest

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/emerishq/emeris-cns-server/cns/auth"
	"github.com/gin-gonic/gin"
)

type AuthHeader struct {
	Token string `header:"Authorization"`
}

func getToken(ctx *gin.Context) (string, error) {
	var token string

	cookie, err := ctx.Request.Cookie("auth._token.google")
	if err != nil {
		a := AuthHeader{}
		if err := ctx.ShouldBindHeader(&a); err != nil {
			return "", errors.New("no auth token found")
		}

		token = a.Token
	} else {
		reqTokenEncoded := cookie.Value

		token, err = url.PathUnescape(reqTokenEncoded)
		if err != nil {
			return "", errors.New("error unescaping token")
		}
	}

	return token, nil
}

func (r *router) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if r.s.a.Env == "test" {
			r.s.l.Infow("skipping auth in test")
			ctx.Set("user", auth.User{
				Name:  "Test Ickle",
				Email: "tester@tendermint.com",
			})
			ctx.Next()
		} else {
			reqToken, err := getToken(ctx)
			if err != nil {
				e(ctx, http.StatusUnauthorized, err)
				r.s.l.Error("token not found", err)
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
			ctx.Set("user", auth.User{
				Name:  name,
				Email: email,
			})

			r.s.l.Infow("incoming request from %s (%s)\n", name, email)
			ctx.Next()
		}
	}
}
