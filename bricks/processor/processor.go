package processor

import (
	"context"
	"sync"

	"github.com/rotisserie/eris"
)

const defaultWorkerCount = 5

type Processor interface {
	Start(ctx context.Context)
	Stop()

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
	p.cancel = cancel

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
}

func (p *processor) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

func (p *processor) Push(i interface{}) {
	p.queue.Push(i)
}

func (p *processor) Done() chan bool {
	return p.doneChan
}

func (p *processor) DoneLength() int64 {
	return p.doneLength
}

func (p *processor) workDoneCallback() {
	p.doneLength++
}

func (p *processor) WorkerOverview() []int64 {
	r := make([]int64, defaultWorkerCount)
	for k := range p.workers {
		r[k] = p.workers[k].DoneLength()
	}
	return r
}
