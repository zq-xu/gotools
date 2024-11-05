package processor

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type TestPicker struct {
}

func (tp *TestPicker) Pickup(i interface{}) {}

// go test -v  -test.run TestProcessor -count=1
func TestProcessor(t *testing.T) {
	Convey("TestProcessor", t, func() {
		p, err := NewProcessor(func() (Picker, error) { return &TestPicker{}, nil })
		So(err, ShouldBeNil)

		go p.Start(context.Background())

		for i := 0; i < 100; i++ {
			p.Push(i)
		}

		time.Sleep(time.Millisecond * 100)
		p.Stop()

		fmt.Println()
		// fmt.Printf("Done length is: %d\n", p.DoneLength())
		ShouldEqual(p.DoneLength(), 100)

		result := p.WorkerOverview()
		for _, v := range result {
			ShouldNotBeZeroValue(v)
		}
	})
}
