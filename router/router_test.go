package router_test

import (
	"testing"

	"github.com/glocurrency/commons/router"
	"github.com/stretchr/testify/assert"
)

func TestRouterWithValidation(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run("routine", func(t *testing.T) {
			t.Parallel()

			r := router.NewRouterWithValidation()
			assert.NotNil(t, r)
		})
	}
}
