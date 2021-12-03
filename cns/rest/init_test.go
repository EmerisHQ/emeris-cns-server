package rest_test

import (
	"fmt"
	"net"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/allinbits/emeris-cns-server/cns/rest"

	"github.com/alicebob/miniredis/v2"
	"github.com/allinbits/emeris-cns-server/cns/chainwatch"
	"github.com/allinbits/emeris-cns-server/cns/config"
	"github.com/allinbits/emeris-cns-server/cns/database"
	"github.com/allinbits/emeris-cns-server/mocks"
	"github.com/allinbits/emeris-cns-server/utils/logging"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// --- Global variables for child tests ---
var testingCtx struct {
	server *rest.Server
}

func TestMain(m *testing.M) {

	// global setup
	server, _, _, tearDown := setup()

	testingCtx.server = server

	// Run test suites
	exitVal := m.Run()

	// clean up
	tearDown()

	os.Exit(exitVal)
}

func setup() (*rest.Server, *gin.Context, *httptest.ResponseRecorder, func()) {

	// --- logger & Gin test context ---
	httpRecorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(httpRecorder)
	logger := logging.New(logging.LoggingConfig{
		LogPath: "",
		Debug:   true,
	})

	// --- CDB ---
	cdbTestServer, err := testserver.NewTestServer()
	checkNoError(err, logger)

	checkNoError(cdbTestServer.WaitForInit(), logger)

	dbConnStr := cdbTestServer.PGURL().String()
	checkNotNil(dbConnStr, "CDB conn. string", logger)

	// connect and run migration
	dbInstance, err := database.New(dbConnStr)
	checkNoError(err, logger)

	// --- Redis ---
	miniRedis, err := miniredis.Run()
	checkNoError(err, logger)
	redisAddr := miniRedis.Addr()

	// --- K8s mock ---
	kube := mocks.Client{}

	// --- Chainwatch process ---
	redisConnection, err := chainwatch.NewConnection(redisAddr)
	checkNoError(err, logger)

	chainwatchInstance := chainwatch.New(
		logger,
		&kube,
		k8sNsInTest,
		redisConnection,
		dbInstance,
		true,
	)

	go chainwatchInstance.Run()

	// --- HTTP server ---
	port, err := getFreePort()
	checkNoError(err, logger)

	conf := &config.Config{
		Debug:                 true,
		DatabaseConnectionURL: dbConnStr,
		KubernetesNamespace:   k8sNsInTest,
		Redis:                 redisAddr,
		LogPath:               "",
		RelayerDebug:          true,
		RESTAddress:           "127.0.0.1:" + port,
	}
	server := rest.NewServer(
		logger,
		dbInstance,
		&kube,
		redisConnection,
		conf,
	)

	ch := make(chan struct{})
	go func() {
		close(ch)
		if err := server.Serve(conf.RESTAddress); err != nil {
			checkNoError(err, logger)
		}
	}()
	<-ch // Wait for the goroutine to start. Still hack!!

	return server, ginCtx, httpRecorder, func() { cdbTestServer.Stop(); miniRedis.Close() }
}

// Empties the DB of data
// Only use in tests executed sequentially
func truncateDB(t *testing.T) {
	_, err := testingCtx.server.DB.Instance.DB.Exec("TRUNCATE cns.chains")
	assert.NoError(t, err)
}

func getFreePort() (port string, err error) {
	ln, err := net.Listen("tcp", ":0")

	if err != nil {
		return "", err
	}

	_, port, _ = net.SplitHostPort(ln.Addr().String())
	_ = ln.Close()

	return port, nil
}

func checkNoError(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}

func checkNotNil(obj interface{}, whatObj string, logger *zap.SugaredLogger) {
	if obj == nil {
		logger.Error(fmt.Printf("Value is nil: %s", whatObj))
		os.Exit(-1)
	}
}
