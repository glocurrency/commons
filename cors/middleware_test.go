package cors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/cors"
	"github.com/stretchr/testify/assert"
)

func TestAllowAllMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	router := gin.New()
	router.Use(cors.AllowAllMiddleware())
	router.GET("/", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
