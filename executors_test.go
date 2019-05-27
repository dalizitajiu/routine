package routine

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/x-mod/errors"
)

func TestExecutor_Execute(t *testing.T) {
	//Guarantee
	assert.Nil(t, Guarantee(ExecutorFunc(func(context.Context) error {
		return nil
	})).Execute(context.TODO()))
	assert.NotNil(t, Guarantee(ExecutorFunc(func(context.Context) error {
		return errors.New("err")
	})).Execute(context.TODO()))
	assert.NotNil(t, Guarantee(ExecutorFunc(func(context.Context) error {
		panic("panic")
		return nil
	})).Execute(context.TODO()))
	assert.NotNil(t, Guarantee(ExecutorFunc(func(context.Context) error {
		panic(errors.New("panic error"))
		return nil
	})).Execute(context.TODO()))

	//Retry
	assert.Equal(t, errors.New("3"), Retry(3, ExecutorFunc(func(ctx context.Context) error {
		return errors.Errorf("%d", FromRetry(ctx))
	})).Execute(context.TODO()))

	//Repeat
	assert.Equal(t, nil, Repeat(3, 10*time.Millisecond, ExecutorFunc(func(ctx context.Context) error {
		Info(ctx, FromRepeat(ctx))
		return nil
	})).Execute(context.TODO()))

	//Concurrent
	assert.Equal(t, nil, Concurrent(3, ExecutorFunc(func(ctx context.Context) error {
		Info(ctx, FromConcurrent(ctx))
		return nil
	})).Execute(context.TODO()))

	//Crontab
	timeoutCtx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
	defer cancel()
	assert.NotNil(t, Crontab("* * * * *", ExecutorFunc(func(ctx context.Context) error {
		return nil
	})).Execute(timeoutCtx))

	//Command
	assert.Equal(t, nil, Command("echo", "hello").Execute(context.TODO()))

	//Timeout
	assert.Nil(t, Timeout(10*time.Millisecond, ExecutorFunc(func(ctx context.Context) error {
		time.Sleep(5 * time.Millisecond)
		return nil
	})).Execute(context.TODO()))
	assert.NotNil(t, Timeout(10*time.Millisecond, ExecutorFunc(func(ctx context.Context) error {
		time.Sleep(15 * time.Millisecond)
		return nil
	})).Execute(context.TODO()))

	//Deadline
	assert.Nil(t, Deadline(time.Now().Add(10*time.Millisecond), ExecutorFunc(func(ctx context.Context) error {
		time.Sleep(5 * time.Millisecond)
		return nil
	})).Execute(context.TODO()))
	assert.NotNil(t, Deadline(time.Now().Add(10*time.Millisecond), ExecutorFunc(func(ctx context.Context) error {
		time.Sleep(15 * time.Millisecond)
		return nil
	})).Execute(context.TODO()))

	//Report
	res1 := make(chan *Result, 1)
	assert.Nil(t, Report(res1, ExecutorFunc(func(context.Context) error {
		return errors.New("err")
	})).Execute(context.TODO()))
	r1 := <-res1
	assert.NotNil(t, r1.Err)

}
