package entity

type CombatResultType int

const (
	CombatResultUnknown CombatResultType = iota
	CombatResultTakeDamage
	CombatResultMessage
)

func (d CombatResultType) String() string {
	return [...]string{"Unknown", "Damage", "Message"}[d]
}

type CombatResult struct {
	Type         CombatResultType
	IntegerValue int32
	StringValue  string
}
