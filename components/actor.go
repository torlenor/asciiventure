package components

// ActionType holds the type of the action to trigger.
type ActionType int

// List of ActionTypes.
const (
	ActionTypeNone ActionType = iota
	ActionTypeMove
	// ActionTypeInteract can either be picking something up, consuming a mutation from the floor or going through a portal
	ActionTypeInteract
	ActionTypeDropItem
	ActionTypeUseItem
)

func (d ActionType) String() string {
	return [...]string{"None", "Move", "Interact", "Drop"}[d]
}

// Actor component tells the systems what action shall be taken next
type Actor struct {
	NextAction ActionType
	IntValue   int
}
