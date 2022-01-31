package auth

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	goauth "google.golang.org/api/oauth2/v2"

	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var conf *oauth2.Config
var secret = []byte("asmiogu;bvzx9vharGDSOJVAG$QY(gadfovzopRASDgfzu^!@^jba90j0awtS{DGa")

var domains = map[string]string{
	"local":   "http://0.0.0.0:8000/",
	"develop": "https://develop--emeris-admin.netlify.app/",
	"staging": "https://staging--emeris-admin.netlify.app/",
	"prod":    "https://admin.emeris.com/",
}

func getRedirectUrl(env string) (string, error) {
	if val, ok := domains[env]; !ok {
		return "", error(fmt.Errorf("invalid environment"))

	} else {
		u, err := url.Parse(val)
		if err != nil {
			return "", err
		}
		u.Path = path.Join(u.Path, "/admin/login")
		s := u.String()

		return s, nil
	}
}

func NewOAuthServer(env string) error {
	url, err := getRedirectUrl(env)

	if err != nil {
		return err
	}

	conf = &oauth2.Config{
		ClientID:     "456830583626-ovlsdesepg4t2g1ufk2nse0b1tbm31pc.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  url,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return nil
}

func Exchange(code string) (*oauth2.Token, error) {
	token, err := conf.Exchange(context.Background(), code, oauth2.AccessTypeOffline)
	return token, err
}

func NewService(token *oauth2.Token) (*goauth.Service, error) {
	svc, err := goauth.NewService(context.Background(), option.WithTokenSource(conf.TokenSource(context.Background(), token)))
	return svc, err
}

func SignJWTs(userInfo *goauth.Userinfo, code string) (string, string, error) {
	authToken := jwt.New(jwt.SigningMethodHS256)

	authClaims := authToken.Claims.(jwt.MapClaims)

	authClaims["sub"] = userInfo.Email
	authClaims["access_uuid"] = code
	authClaims["email"] = userInfo.Email
	authClaims["name"] = userInfo.Name
	authClaims["user"] = userInfo
	authClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	authTokenString, err := authToken.SignedString(secret)

	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)

	refreshClaims := refreshToken.Claims.(jwt.MapClaims)

	refreshClaims["sub"] = userInfo.Email
	refreshClaims["refresh_uuid"] = userInfo.Name
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	refreshTokenString, err := refreshToken.SignedString(secret)

	if err != nil {
		return "", "", err
	}

	return authTokenString, refreshTokenString, nil
}

func ParseJWT(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
