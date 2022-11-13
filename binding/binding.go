package binding

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/erresp"
	"github.com/glocurrency/commons/translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ParseParamUUID parses the param from the request as UUID.
func ParseParamUUID(ctx *gin.Context, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(ctx.Param(name))
	if err != nil {
		return id, fmt.Errorf("cannot parse param: %w", err)
	}

	return id, nil
}

// MustParseParamUUID parses the param from the request as UUID.
// If it fails it will abort the request with an error response.
func MustParseParamUUID(ctx *gin.Context, name string) (uuid.UUID, error) {
	id, err := ParseParamUUID(ctx, name)
	if err != nil {
		ctx.AbortWithStatusJSON(erresp.NewErrResponseBadRequest(fmt.Sprintf("Invalid param %s", name)))
		return id, err
	}

	return id, nil
}

// DecodeBody decodes the request body.
// If it fails it will abort the request with an error response.
func MustDecodeBody(ctx *gin.Context, v interface{}) bool {
	if err := ctx.ShouldBindJSON(v); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := make(map[string]string)

			if translator, ok := translator.GetTranslator(ctx); ok {
				for _, e := range errs {
					formattedErrors[e.Field()] = e.Translate(translator)
				}
			}

			ctx.AbortWithStatusJSON(erresp.NewErrResponseValidationErrors("Request data invalid", formattedErrors))
			return false
		}

		ctx.AbortWithStatusJSON(erresp.NewErrResponseBadRequest("Invalid request body"))
		return false
	}

	return true
}
