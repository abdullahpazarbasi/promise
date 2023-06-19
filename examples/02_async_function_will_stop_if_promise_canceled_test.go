package examples

import (
	"context"
	"fmt"
	"github.com/abdullahpazarbasi/promise"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_Async_Function_Will_Stop_If_Promise_Canceled(t *testing.T) {
	f := promise.New(func(ctx context.Context) (any, error) {
		for i := 1; i < 20; i++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				fmt.Println(">  ", "sub-task tick", i)
				time.Sleep(50 * time.Millisecond)
			}
		}
		fmt.Println(">  ", "this line will not be printed")

		return nil, nil
	})

	fmt.Println(">  ", "sub-task has not started yet")

	p := f.TimeOutLimit(time.Second).Commit()
	fmt.Println(">  ", "sub-task has just started")

	time.Sleep(250 * time.Millisecond)

	p.Cancel()
	fmt.Println(">  ", "promise has just canceled")

	time.Sleep(750 * time.Millisecond)

	_, err := p.Await()
	require.Error(t, err)
}
