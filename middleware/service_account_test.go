package middleware_test

import (
	_ "embed"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	ginfirebasemw "github.com/brokeyourbike/gin-firebase-middleware"
	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/middleware"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/service-account.json
var serviceAccount []byte

func TestRequireIfServiceAccount(t *testing.T) {
	tests := []struct {
		name       string
		header     []byte
		wantId     string
		wantStatus int
		setupMock  func(c *middleware.MockMembersClient)
	}{
		{
			"not service account",
			memberValidId,
			"ab0b166e-c725-4921-b919-fd1cbf43a442",
			http.StatusOK,
			nil,
		},
		{
			"invalid id",
			serviceAccount,
			"1234",
			http.StatusForbidden,
			nil,
		},
		{
			"valid id",
			serviceAccount,
			"ab0b166e-c725-4921-b919-fd1cbf43a442",
			http.StatusOK,
			nil,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			clientMock := middleware.NewMockMembersClient(t)
			if test.setupMock != nil {
				test.setupMock(clientMock)
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// encoding the header value to match what expected by `ginfirebasemw`
			req.Header.Set("X-Apigateway-Api-Userinfo", base64.RawURLEncoding.EncodeToString(test.header))

			w := httptest.NewRecorder()
			router := gin.New()
			router.Use(ginfirebasemw.Middleware())
			router.Use(middleware.RequireIfServiceAccount(test.wantId))
			router.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "the end.")
			})
			router.ServeHTTP(w, req)

			require.Equal(t, test.wantStatus, w.Code)
		})
	}
}
