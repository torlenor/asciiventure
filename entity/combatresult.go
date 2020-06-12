package entity

// CombatResultType is used as enum.
type CombatResultType int

// List of all known CombatResultTypes
const (
	CombatResultUnknown CombatResultType = iota
	CombatResultTakeDamage
	CombatResultMessage
)

func (d CombatResultType) String() string {
	return [...]string{"Unknown", "Damage", "Message"}[d]
}

// CombatResult is one result of a combat action.
type CombatResult struct {
	Type         CombatResultType
	IntegerValue int
	StringValue  string
}
