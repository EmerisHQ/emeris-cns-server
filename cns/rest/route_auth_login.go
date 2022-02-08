package rest

import (
	"errors"
	"fmt"
	"net/http"

	goauth "google.golang.org/api/oauth2/v2"

	"github.com/gin-gonic/gin"
)

type loginResponse struct {
	Token        string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	User         *goauth.Userinfo `json:"user"`
}

func (r *router) Login(ctx *gin.Context) {

	err := ctx.Request.ParseForm()
	if err != nil {
		e(ctx, http.StatusBadRequest, err)
		r.s.l.Error("failed to parse form", err)
		return
	}

	token, err := r.s.a.Exchange(ctx.Request.PostFormValue("code"))
	if err != nil {
		e(ctx, http.StatusBadRequest, err)
		r.s.l.Error("cannot verify code", err)
		return
	}

	oAuth2Service, err := r.s.a.NewService(token)
	if err != nil {
		e(ctx, http.StatusBadRequest, err)
		r.s.l.Error("failed to create oauth service", err)
		return
	}

	userInfo, err := oAuth2Service.Userinfo.Get().Do()
	if err != nil {
		e(ctx, http.StatusBadRequest, err)
		r.s.l.Error("failed to get userinfo", err)
		return
	}

	if userInfo.Hd != "tendermint.com" {
		e(ctx, http.StatusUnauthorized, errors.New("http://youtu.be/otCpCn0l4Wo?t=15"))
		r.s.l.Error("user's domain originates outside tendermint")
		return
	}

	authTokenString, refreshTokenString, err := r.s.a.SignJWTs(userInfo, ctx.Request.PostFormValue("code"))
	if err != nil {
		e(ctx, http.StatusBadRequest, err)
		r.s.l.Error("cannot generate jwts", err)
		return
	}

	ctx.JSON(http.StatusOK, loginResponse{
		fmt.Sprintf("Bearer %s", authTokenString),
		refreshTokenString,
		userInfo,
	})
}
