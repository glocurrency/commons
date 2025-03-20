package audit

import (
	"encoding/json"
)

type BasicEvent struct {
	EventType   string
	ActorType   string
	ActorID     *string
	TargetType  string
	TargetID    *string
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

func NewBasicEventWithUUIDTarget(event string, target UUIDTarget, actorType string, opts ...EventOption) *BasicEvent {
	targetId := target.GetID().String()

	be := &BasicEvent{
		EventType:  event,
		TargetID:   &targetId,
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
	actorID string
}

func (w withActorID) Apply(e *BasicEvent) {
	e.ActorID = &w.actorID
}

func WithActorID(id string) EventOption {
	return withActorID{actorID: id}
}

type withTargetID struct {
	targetID string
}

func (w withTargetID) Apply(e *BasicEvent) {
	e.TargetID = &w.targetID
}

func WithTargetID(id string) EventOption {
	return withTargetID{targetID: id}
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
