package auth

import (
	"context"
	"strings"
	"time"

	goauth "google.golang.org/api/oauth2/v2"

	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type OAServer struct {
	conf   *oauth2.Config
	secret []byte
	Env    string
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewOAuthServer(env, redirectUrl, clientId, clientSecret string, secret []byte) (*OAServer, error) {

	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return &OAServer{
		conf:   conf,
		secret: secret,
		Env:    strings.ToLower(env),
	}, nil
}

func (s *OAServer) Exchange(code string) (*oauth2.Token, error) {
	token, err := s.conf.Exchange(context.Background(), code, oauth2.AccessTypeOffline)
	return token, err
}

func (s *OAServer) NewService(token *oauth2.Token) (*goauth.Service, error) {
	svc, err := goauth.NewService(context.Background(), option.WithTokenSource(s.conf.TokenSource(context.Background(), token)))
	return svc, err
}

func (s *OAServer) SignJWTs(userInfo *goauth.Userinfo, code string) (string, string, error) {
	authToken := jwt.New(jwt.SigningMethodHS256)

	authClaims := authToken.Claims.(jwt.MapClaims)

	authClaims["sub"] = userInfo.Email
	authClaims["access_uuid"] = code
	authClaims["email"] = userInfo.Email
	authClaims["name"] = userInfo.Name
	authClaims["user"] = userInfo
	authClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	authTokenString, err := authToken.SignedString(s.secret)

	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)

	refreshClaims := refreshToken.Claims.(jwt.MapClaims)

	refreshClaims["sub"] = userInfo.Email
	refreshClaims["refresh_uuid"] = userInfo.Name
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	refreshTokenString, err := refreshToken.SignedString(s.secret)
	if err != nil {
		return "", "", err
	}

	return authTokenString, refreshTokenString, nil
}

func (s *OAServer) ParseJWT(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}
