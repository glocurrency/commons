package audit_test

import (
	"testing"

	"github.com/glocurrency/commons/audit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type mockTarget struct{ ID uuid.UUID }

func (m mockTarget) GetID() uuid.UUID {
	return m.ID
}

func (m mockTarget) GetAuditTargetType() string {
	return "mock-target"
}

func TestNewBasicEvent(t *testing.T) {
	target := mockTarget{ID: uuid.New()}

	event := audit.NewBasicEvent(
		"audit-type",
		"target-type",
		"actor-type",
		audit.WithTargetID(uuid.New()),
		audit.WithActorID(uuid.New()),
		audit.WithPayload(target),
		audit.WithPrevPayload(target),
	)
	require.NotNil(t, event)
}

func TestNewBasicEventWithTarget(t *testing.T) {
	target := mockTarget{ID: uuid.New()}

	event := audit.NewBasicEventWithTarget(
		"audit-type",
		target,
		"actor-type",
		audit.WithActorID(uuid.New()),
		audit.WithPayload(target),
		audit.WithPrevPayload(target),
	)
	require.NotNil(t, event)
}
