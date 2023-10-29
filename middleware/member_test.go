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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/member-invalid-id.json
var memberInvalidId []byte

//go:embed testdata/member-valid-id.json
var memberValidId []byte

func TestMemberCtx(t *testing.T) {
	tests := []struct {
		name       string
		header     []byte
		wantStatus int
		setupMock  func(c *middleware.MockMembersClient)
	}{
		{
			"invalid id",
			memberInvalidId,
			http.StatusBadRequest,
			nil,
		},
		{
			"valid id no permissions required",
			memberValidId,
			http.StatusOK,
			func(c *middleware.MockMembersClient) {
				c.On("HasPermissions", mock.Anything, "ten", mock.Anything, mock.Anything).Return(true).Once()
			},
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

			m := middleware.NewMemberCtx("ten", clientMock)
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// encoding the header value to match what expected by `ginfirebasemw`
			req.Header.Set("X-Apigateway-Api-Userinfo", base64.RawURLEncoding.EncodeToString(test.header))

			w := httptest.NewRecorder()
			router := gin.New()
			router.Use(ginfirebasemw.Middleware())
			router.Use(m.Require())
			router.GET("/", func(ctx *gin.Context) {
				id := middleware.MustGetMemberIDFromContext(ctx)
				assert.NotNil(t, id)

				ctx.String(http.StatusOK, "the end.")
			})
			router.ServeHTTP(w, req)

			require.Equal(t, test.wantStatus, w.Code)
		})
	}
}
