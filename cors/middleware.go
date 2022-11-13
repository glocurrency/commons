package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// AllowAllMiddleware is a middleware that allows all origins and headers with default methods.
func AllowAllMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"*"}
	return cors.New(config)
}
