package pathfinding

import (
	"container/heap"

	"github.com/torlenor/asciiventure/utils"
)

func calcCost(current utils.Vec2, next utils.Vec2) float64 {
	// It seems 3/2 as penalty for diagonal movement looks best
	if current.X != next.X && current.Y != next.Y {
		return 3
	}
	return 2
}

// DetermineAstarPath returns the A* path from start to goal.
func DetermineAstarPath(graph Graph, obstacles Obstacles, start utils.Vec2, goal utils.Vec2) []utils.Vec2 {
	if !graph.InDimensions(goal) || graph.Opaque(goal) {
		return []utils.Vec2{}
	}

	open := &positionPriorityQueue{}
	heap.Push(open, &item{value: start, priority: 0})

	cameFrom := map[utils.Vec2]utils.Vec2{}
	costSoFar := map[utils.Vec2]float64{start: 0}

	var current utils.Vec2
	for open.Len() > 0 {
		current = heap.Pop(open).(*item).value.(utils.Vec2)

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
		return []utils.Vec2{}
	}

	if _, ok := cameFrom[goal]; !ok {
		return []utils.Vec2{}
	}

	current = goal
	var path []utils.Vec2
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
