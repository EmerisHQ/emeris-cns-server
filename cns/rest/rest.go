package rest

import (
	"errors"
	"net/http"

	"github.com/emerishq/emeris-cns-server/cns/auth"
	"github.com/emerishq/emeris-cns-server/cns/config"

	"github.com/emerishq/demeris-backend-models/validation"
	"github.com/gin-gonic/gin/binding"

	"github.com/emerishq/emeris-cns-server/cns/chainwatch"

	kube "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/emerishq/emeris-cns-server/cns/database"
	"github.com/emerishq/emeris-utils/logging"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Server struct {
	l          *zap.SugaredLogger
	DB         *database.Instance
	g          *gin.Engine
	KubeClient *kube.Client
	rc         *chainwatch.Connection
	Config     *config.Config
	a          *auth.OAServer
}

type router struct {
	s *Server
}

func NewServer(l *zap.SugaredLogger, d *database.Instance, kube kube.Client, rc *chainwatch.Connection, config *config.Config, a *auth.OAServer) *Server {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.New()

	s := &Server{
		l:          l,
		DB:         d,
		g:          g,
		KubeClient: &kube,
		rc:         rc,
		Config:     config,
		a:          a,
	}

	r := &router{s: s}

	validation.JSONFields(binding.Validator)
	validation.DerivationPath(binding.Validator)
	validation.CosmosRPCURL(binding.Validator)
	validation.Semver(binding.Validator)

	g.Use(logging.LogRequest(l.Desugar()))
	g.Use(ginzap.RecoveryWithZap(l.Desugar(), true))

	// g.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:8000", "http://localhost", "http://127.0.0.1", "http://127.0.0.1:8000"},
	// 	AllowMethods:     []string{"GET", "POST", "OPTIONS"},
	// 	AllowHeaders:     []string{"Authorization"},
	// 	ExposeHeaders:    []string{"Access-Control-Allow-Origin"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	g.Use(func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH,OPTIONS,GET,PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// auth endpoints

	g.GET("/auth/user", r.Auth(), r.User)
	g.POST("/auth/login", r.Login)

	// cns endpoints

	g.GET(r.getChain())
	g.GET(r.getChains())
	g.GET(r.denomsData())
	g.POST(AddChainRoute, r.Auth(), r.addChainHandler)
	g.POST(updatePrimaryChannelRoute, r.Auth(), r.updatePrimaryChannelHandler)
	g.POST(updateDenomsRoute, r.Auth(), r.updateDenomsHandler)
	g.DELETE(DeleteChainRoute, r.Auth(), r.deleteChainHandler)

	g.NoRoute(func(context *gin.Context) {
		e(context, http.StatusNotFound, errors.New("not found"))
	})

	return s
}

func (s *Server) Serve(where string) error {
	return s.g.Run(where)
}

type restError struct {
	Error string `json:"error"`
}

type restValidationError struct {
	ValidationErrors []string `json:"validation_errors"`
}

// e writes err to the caller, with the given HTTP status.
func e(c *gin.Context, status int, err error) {
	var jsonErr interface{}

	jsonErr = restError{
		Error: err.Error(),
	}

	ve := validator.ValidationErrors{}
	if errors.As(err, &ve) {
		rve := restValidationError{}
		for _, v := range ve {
			rve.ValidationErrors = append(rve.ValidationErrors, v.Error())
		}

		jsonErr = rve
	}

	_ = c.Error(err)
	c.AbortWithStatusJSON(status, jsonErr)
}
