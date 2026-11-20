package processor

import (
	"context"
	"sync/atomic"
)

type Worker interface {
	Do(ctx context.Context)

	DoneLength() int64
}

type Picker interface {
	Pickup(interface{})
}

type workerOpt func(*worker)

type worker struct {
	picker Picker

	queue Queue

	callback func()

	doneLength int64
}

func WithWorkerCallback(fn func()) workerOpt {
	return func(w *worker) {
		w.callback = fn
	}
}

func NewWorker(queue Queue, picker Picker, opts ...workerOpt) Worker {
	w := &worker{
		queue:  queue,
		picker: picker,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

func (w *worker) Do(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			i := w.queue.Pop()
			if i == nil {
				// 队列 channel 已关闭，退出
				return
			}

			w.picker.Pickup(i)
			if w.callback != nil {
				w.callback()
			}
			atomic.AddInt64(&w.doneLength, 1)
		}
	}
}

func (w *worker) DoneLength() int64 {
	return atomic.LoadInt64(&w.doneLength)
}
