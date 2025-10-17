package processor

import (
	"context"
	"sync"
)

type Queue interface {
	Start(ctx context.Context)
	Stop()

	Push(interface{})
	Pop() interface{}
}

type queue struct {
	ctx    context.Context
	cancel context.CancelFunc
	closed bool

	lock sync.Mutex
	cond *sync.Cond

	queue []interface{}

	C chan interface{}
}

type queueOpt func(*queue)

func NewQueue(opts ...queueOpt) Queue {
	q := &queue{
		queue: make([]interface{}, 0),
		C:     make(chan interface{}),
	}

	q.cond = sync.NewCond(&q.lock)

	for _, opt := range opts {
		opt(q)
	}

	return q
}

func (q *queue) Start(ctx context.Context) {
	q.ctx, q.cancel = context.WithCancel(ctx)

	for {
		select {
		case <-q.ctx.Done():
			close(q.C)
			return
		default:
			i := q.pop()

			if !q.closed && i != nil {
				q.C <- i
			}
		}
	}
}

func (q *queue) pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	for len(q.queue) == 0 {
		if q.closed {
			return nil
		}

		q.cond.Wait()
	}

	i := q.queue[0]
	q.queue = q.queue[1:]

	return i
}

func (q *queue) Stop() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.cancel()
	q.closed = true
	q.cond.Signal()
}

func (q *queue) Push(i interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queue = append(q.queue, i)
	q.cond.Signal()
}

func (q *queue) Pop() interface{} {
	return <-q.C
}
