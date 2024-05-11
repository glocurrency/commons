package cors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/cors"
	"github.com/glocurrency/commons/router"
	"github.com/stretchr/testify/assert"
)

func TestAllowAllMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	router := router.NewRouterWithValidation()
	router.Use(cors.AllowAllMiddleware())
	router.GET("/", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
