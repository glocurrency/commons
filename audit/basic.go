package audit

import "encoding/json"

type BasicEvent struct {
	EventType   Type
	TargetType  TargetType
	TargetID    string
	ActorType   ActorType
	ActorID     string
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
	actorID string
}

func (w withActorID) Apply(e *BasicEvent) {
	e.ActorID = w.actorID
}

func WithActorID(id string) EventOption {
	return withActorID{actorID: id}
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
