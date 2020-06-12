package components

// Combat holds the properties related to combat.
type Combat struct {
	HP      int `json:"HP"`
	Defense int `json:"Defense"`
	Power   int `json:"Power"`

	CurrentHP int
}
