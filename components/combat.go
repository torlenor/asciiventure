package components

type CombatResultType int

const (
	TakeDamage CombatResultType = iota
	Message
)

func (d CombatResultType) String() string {
	return [...]string{"Damage", "Message"}[d]
}

// Combat holds the properties related to combat.
type Combat struct {
	HP      int32 `json:"HP"`
	Defense int32 `json:"Defense"`
	Power   int32 `json:"Power"`

	CurrentHP int32
}

type CombatResult struct {
	Type         CombatResultType
	IntegerValue int32
	StringValue  string
}
