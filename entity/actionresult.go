package entity

type ActionResultType int

const (
	ActionResultUnknown ActionResultType = iota
	ActionResultItemPickedUp
	ActionResultMutationConsumed
	ActionResultMessage
)

func (d ActionResultType) String() string {
	return [...]string{"Unknown", "ItemPickedUp", "MutationConsumed", "Message"}[d]
}

type ActionResult struct {
	Type         ActionResultType
	IntegerValue int32
	StringValue  string
}
