package game

import (
	"container/heap"
	"math"

	"github.com/torlenor/asciiventure/components"
)

func determineLatticePathSimple(origin components.Position, target components.Position) []components.Position {
	current := origin
	s := []components.Position{}
	for target.X != current.X || target.Y != current.Y {
		if current.X < target.X {
			current.X++
		} else if current.X > target.X {
			current.X--
		}
		if current.Y < target.Y {
			current.Y++
		} else if current.Y > target.Y {
			current.Y--
		}
		lp := components.Position{X: current.X, Y: current.Y}
		s = append(s, lp)
	}
	return s
}

func contains(list []components.Position, p components.Position) bool {
	for _, c := range list {
		if c.X == p.X && c.Y == p.Y {
			return true
		}
	}
	return false
}

func calcCost(current components.Position, next components.Position) float64 {
	// It seems 3/2 as penalty for diagonal movement looks best
	if current.X != next.X && current.Y != next.Y {
		return 3
	}
	return 2
}

func heuristic(a components.Position, b components.Position) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func (g *Game) determineLatticePathAstar(origin components.Position, target components.Position) []components.Position {
	if !g.currentRoom.InDimensions(target) {
		return []components.Position{}
	}

	open := &positionPriorityQueue{}
	heap.Push(open, &item{value: origin, priority: 0})

	cameFrom := map[components.Position]components.Position{}
	costSoFar := map[components.Position]float64{origin: 0}

	var current components.Position
	for open.Len() > 0 {
		current = heap.Pop(open).(*item).value.(components.Position)

		if current.Equal(target) {
			break
		}

		for _, next := range g.currentRoom.Neighbors(current) {
			if g.occupied(next.X, next.Y) && !(next.Equal(target)) {
				continue
			}

			newCost := costSoFar[current] + calcCost(current, next)

			c, inCostSoFar := costSoFar[next]
			if !inCostSoFar || newCost < c {
				costSoFar[next] = newCost
				priority := newCost + heuristic(next, target)
				open.Push(&item{value: next, priority: priority})
				cameFrom[next] = current
			}
		}
	}
	if len(cameFrom) == 0 {
		return []components.Position{}
	}

	if _, ok := cameFrom[target]; !ok {
		return []components.Position{}
	}

	current = target
	var path []components.Position
	for current.X != origin.X || current.Y != origin.Y {
		path = append(path, current)
		current = cameFrom[current]
	}

	// Sort from origin to target
	for i := len(path)/2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]
	}

	return path
}
