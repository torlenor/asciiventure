package game

// MainMenuActionType holds the type of result.
type MainMenuActionType int

// List of MainMenuActionTypes.
const (
	MainMenuActionUnknown MainMenuActionType = iota
	MainMenuActionStartGame
	MainMenuActionLoadGame
	MainMenuActionOptions
	MainMenuActionQuit
)
