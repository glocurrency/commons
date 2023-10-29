package middleware

import (
	"context"

	ginfirebasemw "github.com/brokeyourbike/gin-firebase-middleware"
	"github.com/gin-gonic/gin"
	"github.com/glocurrency/commons/logger"
	"github.com/glocurrency/commons/response"
	"github.com/google/uuid"
)

const MemberIDCtxKey = "memberIDCtx"

type MembersClient interface {
	HasPermissions(ctx context.Context, tenant string, id uuid.UUID, permissions ...[]string) bool
}

type memberCtx struct {
	tenant        string
	membersClient MembersClient
}

func NewMemberCtx(tenant string, membersClient MembersClient) *memberCtx {
	return &memberCtx{tenant: tenant, membersClient: membersClient}
}

func (m *memberCtx) Require(permissions ...[]string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ginfirebasemw.GetUserID(ctx))
		if err != nil {
			logger.WithContext(ctx).
				WithError(err).
				WithField("member_id", ginfirebasemw.GetUserID(ctx)).
				Error("member id cannot be parsed")

			ctx.AbortWithStatusJSON(response.NewErrResponseBadRequest("Member ID not valid"))
			return
		}

		if !m.membersClient.HasPermissions(ctx.Copy(), m.tenant, id, permissions...) {
			ctx.AbortWithStatusJSON(response.NewErrResponseForbidden("Not enough permissions to perform this action"))
			return
		}

		ctx.Set(MemberIDCtxKey, id)
		ctx.Next()
	}
}

// MustGetAdminFromContext returns the member ID from the context.
func MustGetMemberIDFromContext(ctx *gin.Context) uuid.UUID {
	return ctx.MustGet(MemberIDCtxKey).(uuid.UUID)
}
