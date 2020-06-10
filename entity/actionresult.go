package entity

import "github.com/torlenor/asciiventure/components"

// ActionResultType holds the type of result.
type ActionResultType int

// List of ActionResultTypes.
const (
	ActionResultUnknown ActionResultType = iota
	ActionResultItemPickedUp
	ActionResultMutationConsumed
	ActionResultMessage
)

func (d ActionResultType) String() string {
	return [...]string{"Unknown", "ItemPickedUp", "MutationConsumed", "Message"}[d]
}

// ActionResult is the result of an action.
// It can have a value.
type ActionResult struct {
	Type                ActionResultType
	MutationEffectValue components.MutationEffect
	IntegerValue        int32
	StringValue         string
}
