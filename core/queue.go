package core

type queue struct {
	data []int
}

func newQueue(data []int) *queue {
	return &queue{data: data}
}

func (q *queue) addFirst(n int) {
	q.data = append([]int{n}, q.data...)
}

func (q *queue) addLast(n int) {
	q.data = append(q.data, n)
}

func (q *queue) popNonZero() (int, bool) {
	length := len(q.data)
	for i := 0; i < length; i++ {
		if q.data[i] != 0 {
			n := q.data[i]
			q.data = q.data[i+1:]
			return n, true
		}
	}
	return -1, false
}
