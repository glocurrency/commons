package translator

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

// translatorCtx stores the translator in the context
const translatorCtx = "translatorCtx"

// SetTranslatorMiddleware is a middleware that adds translator to the context.
func SetTranslatorMiddleware(t ut.Translator) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(translatorCtx, t)
		c.Next()
	}
}

// GetTranslator returns the translator from the context.
func GetTranslator(ctx *gin.Context) (ut.Translator, bool) {
	value, ok := ctx.Get(translatorCtx)
	if !ok {
		return nil, false
	}

	t, ok := value.(ut.Translator)
	if !ok {
		return nil, false
	}

	return t, true
}
