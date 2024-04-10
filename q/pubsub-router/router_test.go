package pubsubrouter_test

import (
	"bytes"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	pubsubrouter "github.com/glocurrency/commons/q/pubsub-router"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)
	os.Exit(m.Run())
}

//go:embed testdata/with-name-d.json
var withNameD []byte

//go:embed testdata/with-name-e.json
var withNameE []byte

//go:embed testdata/with-name-def.json
var withNameDEF []byte

func TestRouting(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		body       []byte
		wantStatus int
	}{
		{"pubsub not found", "/abc", nil, http.StatusNotFound},
		{"pubsub not found, name do not match route", "/abc", withNameE, http.StatusNotFound},
		{"pubsub routed", "/abc", withNameD, http.StatusAccepted},
		{"pubsub routed", "/abc?a=b", withNameD, http.StatusAccepted},
		{"pubsub routed", "/abc/", withNameD, http.StatusAccepted},
		{"pubsub routed", "/abc/?a=b", withNameD, http.StatusAccepted},
		{"pubsub routed", "/abc", withNameDEF, http.StatusAccepted},
		{"not pubpub", "/123/4", withNameD, http.StatusOK},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, test.url, bytes.NewReader(test.body))

			router := pubsubrouter.NewRouter(gin.New(), "/abc")

			group1 := router.Group("/abc")
			group1.POST("d", func(ctx *gin.Context) { ctx.Status(http.StatusAccepted) })
			group1.POST("d-e-f", func(ctx *gin.Context) { ctx.Status(http.StatusAccepted) })

			group2 := router.Group("/123")
			group2.POST("4", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

			router.ServeHTTP(w, req)
			assert.Equal(t, test.wantStatus, w.Code)
		})
	}
}
