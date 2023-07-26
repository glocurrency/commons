package audit

import "github.com/google/uuid"

type Type string
type TargetType string
type ActorType string

type Target interface {
	GetID() uuid.UUID
	GetAuditTargetType() TargetType
}
