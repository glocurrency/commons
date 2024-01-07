package audit

import "github.com/google/uuid"

const (
	ActorTypeServer  = "SERVER"
	ActorTypeUser    = "USER"
	ActorTypeMember  = "MEMBER"
	ActorTypeWebhook = "WEBHOOK"
)

type Target interface {
	GetID() uuid.UUID
	GetAuditTargetType() string
}
