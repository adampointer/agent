package support

import (
	"bytes"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func GetTestDataReadClose(t testing.TB, filename string) io.ReadCloser {
	t.Helper()

	data := GetTestData(t, filename)
	buf := bytes.NewBuffer(data)
	return io.NopCloser(buf)
}

func GetTestData(t testing.TB, filename string) []byte {
	t.Helper()

	wd, err := os.Getwd()
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(wd, "..", "..", "..", "..", "..", "testdata", filename))
	require.NoError(t, err)
	return data
}

func TestServer(t testing.TB, path, filename string) *httptest.Server {
	t.Helper()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET(path, func(ginCtx *gin.Context) {
		_, err := ginCtx.Writer.Write(GetTestData(t, filename))
		require.NoError(t, err)
	})
	srv := httptest.NewServer(r)

	t.Cleanup(func() {
		srv.Close()
	})
	return srv

}
