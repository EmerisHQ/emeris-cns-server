package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type AuthHeader struct {
	Token string `header:"Authorization"`
}

type AuthClient interface {
	AuthUser() gin.HandlerFunc
}

type Auth struct {
	Client auth.Client
}

func NewAuthClient() (AuthClient, error) {

	a := Auth{}
	opt := option.WithCredentialsJSON([]byte(`{
		"type": "service_account",
		"project_id": "emeris-admin-ui",
		"private_key_id": "e7cd8fdf33e2a1a76f30f47c6718b21938bf6f23",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCh4cF7h2ExyVy5\nXZZyujKzjM8m/8CQYiwYD90rBrsN0ixZ9L2i62gnzxG7MSkX2JQHl9mzfnwoREhu\n/UYbPO0n0RJsoZ8ZdoIDpGVMRi21Srwf2+fkAEv+x2GJGxOTFl+QdCbzefwBW/Ss\nh/yAEJ8BgtGyGBDHu7OzXlyrgn+rItsJUcldaGxa4lpgiWV/NEccaYAvDfhWgKmK\n9dp7UksMDEWIdBDPwTH+nQ6C+LWdS4beq4KtYzdhcAXflZWb4VVj3cZjlA4uc8JA\nYSLGPW6LAwu2O0ugQIqhIBBeadydkCZK4kynUnA2APw7A4uqL+c2wAN02ctJI+Sy\nWequFDxHAgMBAAECggEAGwcURMmfoq5Z+uDzQ4hu+qdh1sMQpYqejg3oAU0IYhBb\nM1G3b8IaC7t43GYi1EZmwLXLtTpDBH4SEeXblKShe+peRyDc7WVp463I8+krrH8j\n1bXji5+5EHq9gCSzKfWsUvPxpOkS+C8gNMYnlEIyKhBrbm6yLobaQ/JXSpNpOWs+\nXJRyklVGrtFvik0v6XFnwuULRHE4g0+g/732yV79Y35KhJet9hvMpy8Y0/ktkC3Y\nxu+00rfeVQBMVT52bJEyUKCHF253e5IcvMggzMP9TqwLglwLtPKsVyLx0rWF60Ro\nzUa8LSd8JprhA9iSuuqz+mR3riGcXl1wrkrZQKOVQQKBgQDTRp0fjFdnidGi+lBy\nSIakqpbkj9yjNT1wbWkeuV1J1ASPku4jfprWHETBEaiZigSHzGCLFt9Qg5Bch6Fw\nkk1iI2ubVrl87ysNMvsGXGPwqP5V7i+nFpXMYM3z7hY1TDfbOB8HORL5flUuqwXv\nzI4R7fCEfLQyaBiu0jpwpWoSEQKBgQDEJmaiQ83xHW1nQB3aNjnvpQqRSvSS7hGR\nwSjR2+6TktB+6SpXfk9Uv62CwrsuVtvp2zUG+P8gum/v9rNwWUYdOnvQXjBc4skN\nDCl6ZeaLh91yfBy4ZrROMD2v75HrVchP5nYfrAWBJmdy0MK7FotyDzLtenrLRlpz\nsbdT9P4Q1wKBgClyq/Z5eNg2IGthwhB5i/iYAtw6IOXf1vrMbBf784I9Vtu3zoIm\nH0gr6Y0a4sGkYvklLjd7ODo6ZULR1OkZuparLjweSmtpHEANpVN9Ipoe/S5seOrF\nsoOS5jSZm7+/ASI/o06ucruBfkKWiKafsatwy4OiV1OgOl9pnM9mlCWRAoGAeMQj\n4Lfabh9uImnpd1Z3qUJ2FSqPFn+ZNaI1na/JXfbAg8LPHPtZoJY7IA0A7fDwiTU7\nmsVnXyEqlhXQONXeQ1Skso+rOyUuH+hjCUcAANxvzXL4w9gIHzO4Z0AbGUfBguAj\nzjA9W1znyFsb6dBhnqIY+vmz7L+uJRlABGMMohUCgYAwVEm4agIuvNRkhFvcQMrQ\nQfZDUXpk8X7XXMQ+dMxH9tEGc0y1SRWPl7K1657ZtAFfDWb0zvszrmhofZUf6Gq2\nH4mXOYZUx+ZWNWlQUxOk15qP3EcTDPXNmOD3HPoc1B1IPrwwLJciO3sitND4xl7m\nrm7tWPEjqPQ5m+/gl8IemQ==\n-----END PRIVATE KEY-----\n",
		"client_email": "firebase-adminsdk-eg4xh@emeris-admin-ui.iam.gserviceaccount.com",
		"client_id": "102725695233406839451",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-eg4xh%40emeris-admin-ui.iam.gserviceaccount.com"
	  }
	  `))
	fire, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return a, err
	}

	client, err := fire.Auth(context.Background())

	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
		return a, err
	}

	a.Client = *client

	return a, nil
}

func (a Auth) AuthUser() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		auth := AuthHeader{}

		if err := ctx.ShouldBindHeader(&auth); err != nil {
			e(ctx, http.StatusUnauthorized, err)
			return
		}

		jwtTokenHeader := strings.Split(auth.Token, "JWT ")

		token, err := a.Client.VerifyIDToken(ctx, jwtTokenHeader[1])

		if err != nil {
			e(ctx, http.StatusUnauthorized, err)
			return
		}

		ctx.Set("user", token)
		ctx.Next()
	}
}

type AuthError struct {
	Error string `json:"error"`
}

func e(c *gin.Context, status int, err error) {
	jsonErr := AuthError{
		Error: err.Error(),
	}

	_ = c.Error(err)
	c.AbortWithStatusJSON(status, jsonErr)
}
