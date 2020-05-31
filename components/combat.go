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

func (c *Combat) Attack(target *Combat) (results []CombatResult) {
	if target == nil {
		return
	}

	dmg := c.Power - target.Defense
	if dmg < 0 {
		dmg = 0
	}

	results = append(results, CombatResult{Type: TakeDamage, IntegerValue: dmg})

	return
}
