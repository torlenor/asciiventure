package components

type CombatResultType int

const (
	TakeDamage CombatResultType = iota
	Message
)

func (d CombatResultType) String() string {
	return [...]string{"Damage", "Message"}[d]
}

type Combat struct {
	MaxHP   int32
	HP      int32
	Defense int32
	Power   int32
}

type CombatResult struct {
	Type         CombatResultType
	IntegerValue int32
	StringValue  string
}
