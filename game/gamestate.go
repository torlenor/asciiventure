package game

type gameState int

const (
	mainMenu gameState = iota
	playersTurn
	enemyTurn
	gameOver
)

func (d gameState) String() string {
	return [...]string{"mainMenu", "playersTurn", "enemyTurn", "gameOver"}[d]
}
