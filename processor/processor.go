package processor

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/rotisserie/eris"
)

const defaultWorkerCount = 5

type Processor interface {
	Start(ctx context.Context)
	Stop()

	// Done 返回的 channel 在所有 worker 退出（Stop / ctx 取消之后）时收到信号。
	Done() chan bool

	Push(interface{})
	DoneLength() int64

	WorkerOverview() []int64
}

type PickerGenerateFn func() (Picker, error)
type processorOpt func(*processor)

type processor struct {
	queue Queue

	workers []Worker

	wg     sync.WaitGroup
	mu     sync.Mutex
	cancel context.CancelFunc

	pickerGenerator PickerGenerateFn

	doneLength int64
	doneChan   chan bool
}

func NewProcessor(fn PickerGenerateFn, opts ...processorOpt) (Processor, error) {
	p := &processor{
		pickerGenerator: fn,
		workers:         make([]Worker, defaultWorkerCount),
		queue:           NewQueue(),
		doneChan:        make(chan bool, 1),
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.pickerGenerator == nil {
		return nil, eris.New("empty pickerGenerator")
	}

	for k := range p.workers {
		picker, err := p.pickerGenerator()
		if err != nil {
			return nil, eris.Wrap(err, "failed to generate picker")
		}

		p.workers[k] = NewWorker(p.queue, picker, WithWorkerCallback(p.workDoneCallback))
	}

	return p, nil
}

func (p *processor) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	p.mu.Lock()
	p.cancel = cancel
	p.mu.Unlock()

	go p.queue.Start(ctx)
	defer p.queue.Stop()

	for index := range p.workers {
		p.wg.Add(1)

		go func(i int) {
			defer p.wg.Done()
			p.workers[i].Do(ctx)
		}(index)
	}

	p.wg.Wait()
	p.signalDone()
}

// signalDone 通知等待者处理器已完全停止；重复 Start 时只通知一次。
func (p *processor) signalDone() {
	select {
	case p.doneChan <- true:
	default:
	}
}

func (p *processor) Stop() {
	p.mu.Lock()
	cancel := p.cancel
	p.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	// 唤醒阻塞在 pop/Pop 上的 queue 与 worker，否则它们收不到 ctx 取消
	p.queue.Stop()
}

func (p *processor) Push(i interface{}) {
	p.queue.Push(i)
}

func (p *processor) Done() chan bool {
	return p.doneChan
}

func (p *processor) DoneLength() int64 {
	return atomic.LoadInt64(&p.doneLength)
}

func (p *processor) workDoneCallback() {
	atomic.AddInt64(&p.doneLength, 1)
}

func (p *processor) WorkerOverview() []int64 {
	r := make([]int64, defaultWorkerCount)
	for k := range p.workers {
		r[k] = p.workers[k].DoneLength()
	}
	return r
}
