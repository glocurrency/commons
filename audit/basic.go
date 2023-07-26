package audit

import (
	"encoding/json"

	"github.com/google/uuid"
)

type BasicEvent struct {
	EventType   Type
	TargetType  TargetType
	TargetID    uuid.UUID
	ActorType   ActorType
	ActorID     uuid.NullUUID
	PrevPayload json.RawMessage
	Payload     json.RawMessage
}

func NewBasicEvent(event Type, target Target, actor ActorType, opts ...EventOption) *BasicEvent {
	be := &BasicEvent{
		EventType:  event,
		TargetID:   target.GetID(),
		TargetType: target.GetAuditTargetType(),
		ActorType:  actor,
	}

	for _, o := range opts {
		o.Apply(be)
	}

	return be
}

type EventOption interface {
	Apply(*BasicEvent)
}

type withActorID struct {
	actorID uuid.NullUUID
}

func (w withActorID) Apply(e *BasicEvent) {
	e.ActorID = w.actorID
}

func WithActorID(id uuid.UUID) EventOption {
	return withActorID{actorID: uuid.NullUUID{UUID: id, Valid: true}}
}

type withPrevPayload struct{ payload json.RawMessage }

func (w withPrevPayload) Apply(e *BasicEvent) {
	e.PrevPayload = w.payload
}

func WithPrevPayload(v interface{}) EventOption {
	return withPrevPayload{payload: TryMarshall(v)}
}

type withPayload struct{ payload json.RawMessage }

func (w withPayload) Apply(e *BasicEvent) {
	e.Payload = w.payload
}

func WithPayload(v interface{}) EventOption {
	return withPayload{payload: TryMarshall(v)}
}
