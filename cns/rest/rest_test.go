package rest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

const k8sNsInTest = "emeris"

func TestServerSetup(t *testing.T) {
	t.Parallel()

	// act
	resp, err := http.Get(fmt.Sprintf("http://%s%s", testingCtx.server.config.RESTAddress, "/foo/bar"))

	// assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}



