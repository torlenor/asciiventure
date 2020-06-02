package game

// An Item is something we manage in a priority queue.
type item struct {
	value    interface{} // The value of the item; arbitrary.
	priority float64     // The priority of the item in the queue.
}

type positionPriorityQueue []*item

func (pq positionPriorityQueue) Len() int { return len(pq) }

func (pq positionPriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq positionPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *positionPriorityQueue) Push(x interface{}) {
	item := x.(*item)
	*pq = append(*pq, item)
}

func (pq *positionPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
