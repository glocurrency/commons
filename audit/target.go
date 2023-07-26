package audit

type Type string
type TargetType string
type ActorType string

type Target interface {
	GetID() string
	GetAuditTargetType() TargetType
}
