package binding

import (
	"fmt"

	"github.com/gin-gonic/gin"
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
