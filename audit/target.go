package audit

import "github.com/google/uuid"

type Target interface {
	GetID() uuid.UUID
	GetAuditTargetType() string
}
