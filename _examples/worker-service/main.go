package main

import (
	"context"
	"time"

	"github.com/zq-xu/gotools/logx"
	"github.com/zq-xu/gotools/processor"
	"github.com/zq-xu/gotools/utilsx"
)

type Task struct {
	ID   int64
	Name string
}

type taskPicker struct {
	// 真实场景里这里通常持有 DB 连接、HTTP client 之类
}

func (p *taskPicker) Pickup(i interface{}) {
	task, ok := i.(*Task)
	if !ok {
		logx.Logger.Warnf("unexpected job type: %T", i)
		return
	}

	logx.Logger.Infof("handle task %d/%s", task.ID, task.Name)
	time.Sleep(50 * time.Millisecond)
}

func main() {
	p, err := processor.NewProcessor(func() (processor.Picker, error) {
		return &taskPicker{}, nil
	})
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	utilsx.Notify(cancel)

	// 生产者在别的 goroutine 里持续投递任务
	go func() {
		for i := 0; i < 20; i++ {
			p.Push(&Task{ID: int64(i), Name: "demo"})
		}
	}()

	// Start 会阻塞到 ctx 被取消，内部是 5 个固定 worker 消费同一个队列
	p.Start(ctx)

	logx.Logger.Infof("done=%d perWorker=%v", p.DoneLength(), p.WorkerOverview())
}
