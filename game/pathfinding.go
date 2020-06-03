package game

import (
	"container/heap"

	"github.com/torlenor/asciiventure/components"
)

type Graph interface {
	InDimensions(p components.Position) bool
	Neighbors(p components.Position) []components.Position
	Distance(a components.Position, b components.Position) float64
}

type Obstacles interface {
	Occupied(p components.Position) bool
}

func determineStraightLinePath(start components.Position, goal components.Position) []components.Position {
	current := start
	s := []components.Position{}
	for goal.X != current.X || goal.Y != current.Y {
		if current.X < goal.X {
			current.X++
		} else if current.X > goal.X {
			current.X--
		}
		if current.Y < goal.Y {
			current.Y++
		} else if current.Y > goal.Y {
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

func determineAstarPath(graph Graph, obstacles Obstacles, start components.Position, goal components.Position) []components.Position {
	if !graph.InDimensions(goal) {
		return []components.Position{}
	}

	open := &positionPriorityQueue{}
	heap.Push(open, &item{value: start, priority: 0})

	cameFrom := map[components.Position]components.Position{}
	costSoFar := map[components.Position]float64{start: 0}

	var current components.Position
	for open.Len() > 0 {
		current = heap.Pop(open).(*item).value.(components.Position)

		if current.Equal(goal) {
			break
		}

		for _, next := range graph.Neighbors(current) {
			if obstacles.Occupied(next) && !(next.Equal(goal)) {
				continue
			}

			newCost := costSoFar[current] + calcCost(current, next)

			c, inCostSoFar := costSoFar[next]
			if !inCostSoFar || newCost < c {
				costSoFar[next] = newCost
				priority := newCost + graph.Distance(next, goal)
				open.Push(&item{value: next, priority: priority})
				cameFrom[next] = current
			}
		}
	}
	if len(cameFrom) == 0 {
		return []components.Position{}
	}

	if _, ok := cameFrom[goal]; !ok {
		return []components.Position{}
	}

	current = goal
	var path []components.Position
	for current.X != start.X || current.Y != start.Y {
		path = append(path, current)
		current = cameFrom[current]
	}

	// Sort from start to goal
	for i := len(path)/2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]
	}

	return path
}
