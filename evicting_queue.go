package main

type EvictingQueue struct {
	startIndex int
	endIndex   int
	data       []string
	maxSize    int
}

func NewEvictingQueue(size int) EvictingQueue {
	return EvictingQueue{
		startIndex: 0,
		endIndex:   0,
		data:       make([]string, size, size),
		maxSize:    size,
	}
}

func (q *EvictingQueue) Enqueue(value string) {
	q.startIndex = (q.startIndex + q.maxSize - 1) % q.maxSize
	if q.endIndex == q.startIndex {
		q.endIndex = (q.endIndex + q.maxSize - 1) % q.maxSize
	}
	q.data[q.startIndex] = value
}

func (q *EvictingQueue) Pop() string {
	if q.Length() == 0 {
		panic("Cannot pop from an empty queue")
	}
	value := q.GetMostRecentItem(0)
	q.startIndex = (q.startIndex + 1) % q.maxSize
	return value
}

func (q *EvictingQueue) GetMostRecentItem(index int) string {
	actualIndex := (index + q.startIndex) % q.maxSize
	return q.data[actualIndex]
}

func (q *EvictingQueue) Length() int {
	if q.endIndex >= q.startIndex {
		return q.endIndex - q.startIndex
	}
	return (q.maxSize - q.startIndex) + q.endIndex
}

func (q *EvictingQueue) Iter() []string {
	/*
		TODO this could be made better without the malloc by passing ina  callback based on its
		usage.  But, I'm not super confident without some testing how to handle callbacks
		intended to break out of the loop
	*/
	returnValues := make([]string, 0)
	for i := 0; i < q.Length(); i++ {
		returnValues = append(
			returnValues,
			q.GetMostRecentItem(i),
		)
	}
	return returnValues
}
