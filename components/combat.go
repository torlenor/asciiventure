package components

// Combat holds the properties related to combat.
type Combat struct {
	HP      int32 `json:"HP"`
	Defense int32 `json:"Defense"`
	Power   int32 `json:"Power"`

	CurrentHP int32
}
