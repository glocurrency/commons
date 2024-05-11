package router_test

import (
	"testing"

	"github.com/glocurrency/commons/router"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestRouterWithValidation(t *testing.T) {
	eg := errgroup.Group{}
	for i := 0; i < 100; i++ {
		eg.Go(func() error {
			r := router.NewRouterWithValidation()
			assert.NotNil(t, r)
			return nil
		})
	}
	assert.NoError(t, eg.Wait())
}
