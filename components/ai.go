package components

// The AI component holds information which influences the AI system for a given entity.
type AI struct {
	AttackRange      int32 `json:"AttackRange"`
	AttackRangeUntil int32 `json:"AttackRangeUntil"`
}
