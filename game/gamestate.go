package game

type gameState int

const (
	playersTurn gameState = iota
	enemyTurn
	gameOver
)

func (d gameState) String() string {
	return [...]string{"playersTurn", "enemyTurn", "gameOver"}[d]
}
