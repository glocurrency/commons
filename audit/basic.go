package audit

import (
	"encoding/json"

	"github.com/google/uuid"
)

type BasicEvent struct {
	EventType   string
	ActorType   string
	ActorID     uuid.NullUUID
	TargetType  string
	TargetID    uuid.NullUUID
	PrevPayload json.RawMessage
	Payload     json.RawMessage
}

func NewBasicEvent(event string, targetType string, actorType string, opts ...EventOption) *BasicEvent {
	be := &BasicEvent{
		EventType:  event,
		TargetType: targetType,
		ActorType:  actorType,
	}

	for _, o := range opts {
		o.Apply(be)
	}

	return be
}

func NewBasicEventWithTarget(event string, target Target, actorType string, opts ...EventOption) *BasicEvent {
	be := &BasicEvent{
		EventType:  event,
		TargetID:   uuid.NullUUID{UUID: target.GetID(), Valid: true},
		TargetType: target.GetAuditTargetType(),
		ActorType:  actorType,
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

type withTargetID struct {
	targetID uuid.NullUUID
}

func (w withTargetID) Apply(e *BasicEvent) {
	e.TargetID = w.targetID
}

func WithTargetID(id uuid.UUID) EventOption {
	return withTargetID{targetID: uuid.NullUUID{UUID: id, Valid: true}}
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
