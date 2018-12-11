package util

type FIFOQueueError struct {
	msg string
}

func (e *FIFOQueueError) Error() string {
	return e.msg
}

type FIFOQueue struct {
	queue []interface{}
}

func NewFIFOQueue() FIFOQueue {
	return FIFOQueue{}
}

func (q *FIFOQueue) Size() int {
	if q.queue == nil {
		return 0
	}

	return len(q.queue)
}

func (q *FIFOQueue) Pop(num int) ([]interface{}, error) {
	if num < 1 {
		return nil, &FIFOQueueError{
			msg: "num should be greater than 0",
		}
	}

	if q.queue == nil {
		return nil, &FIFOQueueError{
			msg: "queue is empty",
		}
	}

	if num > len(q.queue) {
		num = len(q.queue)
	}

	qSlice := q.queue[:num]
	q.queue = q.queue[num+1:]
	return qSlice, nil
}

func (q *FIFOQueue) Push(ele interface{}) error {
	if q.queue != nil {
		q.queue = append(q.queue, ele)
		return nil
	}

	q.queue = []interface{}{ele}
	return nil
}
