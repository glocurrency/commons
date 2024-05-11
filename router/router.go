package router

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/glocurrency/commons/validation"
	"github.com/go-playground/validator/v10"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("alphanumspace", validation.ValidateAlphaNumSpace)
		v.RegisterValidation("alphanumspacedash", validation.ValidateAlphaNumSpaceDash)
		v.RegisterValidation("banksupported", validation.ValidateBankSupported)
	}
}

func NewRouterWithValidation() *gin.Engine {
	return sync.OnceValue(func() *gin.Engine { return gin.New() })()
}
