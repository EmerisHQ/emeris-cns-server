package rest_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const k8sNsInTest = "emeris"

func TestServerSetup(t *testing.T) {
	t.Parallel()

	// act
	resp, err := http.Get(fmt.Sprintf("http://%s%s", testingCtx.server.Config.RESTAddress, "/foo/bar"))

	// assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
