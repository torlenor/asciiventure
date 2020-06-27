package components

// Health holds the properties related to health of an entity.
type Health struct {
	HP           int32 `json:"HP"`
	CurrentHP    int32 `json:"CurrentHP"`
	Regeneration int32 `json:"Regeneration"`
}
