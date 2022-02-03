package config

import (
	"fmt"

	"github.com/allinbits/demeris-backend-models/validation"
	"github.com/allinbits/emeris-utils/configuration"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	DatabaseConnectionURL string `validate:"required"`
	Redis                 string `validate:"required"`
	KubernetesNamespace   string `validate:"required"`
	LogPath               string
	RESTAddress           string `validate:"required"`
	Debug                 bool
	RelayerDebug          bool
	Env                   string `validate:"required"`
	Secret                string `validate:"required"`
	RedirectURL           string `validate:"required"`
	OAuth2ClientID        string `validate:"required"`
	OAuth2ClientSecret    string `validate:"required"`
}

func (c Config) Validate() error {
	err := validator.New().Struct(c)
	if err == nil {
		return nil
	}
	return fmt.Errorf(
		"configuration file error: %w",
		validation.MissingFieldsErr(err, false),
	)
}

func ReadConfig() (*Config, error) {
	var c Config

	return &c, configuration.ReadConfig(&c, "demeris-cns", map[string]string{
		"RESTAddress":         ":9999",
		"KubernetesNamespace": "emeris",
		"RelayerDebug":        "true",
		"Env":                 "local",
		"Secret":              "asmiogu;bvzx9vharGDSOJVAG$QY(gadfovzopRASDgfzu^!@^jba90j0awtS{DGa",
	})
}
