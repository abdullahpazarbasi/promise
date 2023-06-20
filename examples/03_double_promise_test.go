package examples

import (
	"context"
	"fmt"
	"github.com/abdullahpazarbasi/promise/v3"
	"testing"
	"time"
)

func Test_Double_Promise(t *testing.T) {
	f1 := promise.New(func(ctx context.Context) (string, error) {
		// step 1
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			time.Sleep(100 * time.Millisecond)
		}
		// step 2
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			time.Sleep(100 * time.Millisecond)
		}

		return "OK", nil
	})
	f2 := promise.New(func(ctx context.Context) (bool, error) {
		// step 1
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			time.Sleep(100 * time.Millisecond)
		}
		// step 2
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			time.Sleep(100 * time.Millisecond)
		}

		return true, nil
	})
	p1 := f1.TimeOutLimit(500 * time.Millisecond).Commit()
	p2 := f2.TimeOutLimit(400 * time.Millisecond).Commit()

	fmt.Println("Doing something on primary parallel path ...")
	time.Sleep(200 * time.Millisecond)

	out1, err1 := p1.Await()
	if err1 != nil {
		return
	}
	fmt.Println("Waiting for the other committed promise, may be the task is already done a long time ago")
	out2, err2 := p2.Await()
	if err2 != nil {
		return
	}
	fmt.Printf("Output of async function 1: %v\n", out1)
	fmt.Printf("Output of async function 2: %v\n", out2)
}
