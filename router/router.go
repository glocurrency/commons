package router

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/glocurrency/commons/translator"
	"github.com/glocurrency/commons/validation"
	"github.com/go-playground/validator/v10"
)

func NewRouterWithValidation() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	sync.OnceFunc(func() {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			v.RegisterValidation("alphanumspace", validation.ValidateAlphaNumSpace)
			v.RegisterValidation("alphanumspacedash", validation.ValidateAlphaNumSpaceDash)
			v.RegisterValidation("banksupported", validation.ValidateBankSupported)

			t := translator.RegisterTranslatorFor(v)
			router.Use(translator.SetTranslatorMiddleware(t))
		}
	})
	return router
}
