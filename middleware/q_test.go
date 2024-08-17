package middleware_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/middleware"
	"github.com/glocurrency/commons/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestQCtx(t *testing.T) {
	tests := []struct {
		name       string
		body       []byte
		wantStatus int
		setupMock  func(lockerMock *middleware.MockLocker)
	}{
		{
			"success without lock",
			validPubSubMsg,
			http.StatusOK,
			nil,
		},
		{
			"success with lock",
			validPubSubUniqueMsg,
			http.StatusOK,
			func(lockerMock *middleware.MockLocker) {
				lockerMock.On("TryToLock", mock.Anything, "unique-456").Return(nil)
			},
		},
		{
			"invalid message",
			invalidPubSubMsg,
			http.StatusBadRequest,
			nil,
		},
		{
			"cannot lock by unique the key",
			validPubSubUniqueMsg,
			http.StatusUnprocessableEntity,
			func(lockerMock *middleware.MockLocker) {
				lockerMock.On("TryToLock", mock.Anything, "unique-456").Return(errors.New("I am an error"))
			},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			tasksMock := middleware.NewMockLocker(t)
			if test.setupMock != nil {
				test.setupMock(tasksMock)
			}

			m := middleware.NewQCtx(tasksMock)
			req := httptest.NewRequest(http.MethodPost, "/", nil)

			if test.body != nil {
				req.Body = io.NopCloser(bytes.NewReader(test.body))
			}

			w := httptest.NewRecorder()
			router := router.NewRouterWithValidation()
			router.Use(m.RequireValidMessage())
			router.POST("/", func(ctx *gin.Context) {
				msg := middleware.MustGetQMessageFromContext(ctx)
				assert.NotNil(t, msg)

				ctx.String(http.StatusOK, "the end.")
			})
			router.ServeHTTP(w, req)

			require.Equal(t, test.wantStatus, w.Code)
		})
	}
}
