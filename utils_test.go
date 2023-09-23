package promise

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergeAllResolved(t *testing.T) {
	promises := make([]*Promise[int], 0, 10)

	for i := 0; i < 10; i++ {
		promise := New(func(resolve func(int), reject func(error)) {
			time.Sleep(time.Millisecond * 100)
			resolve(rand.Intn(10))
		})
		promises = append(promises, promise)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	values, err := MergeAll(ctx, promises...).Await(ctx)
	assert.NotEmpty(t, values)
	assert.Nil(t, err)
}

func TestMergeAllWithOneRejected(t *testing.T) {
	promises := make([]*Promise[int], 0, 10)

	for i := 0; i < 9; i++ {
		promise := New(func(resolve func(int), reject func(error)) {
			time.Sleep(time.Millisecond * 100)
			resolve(rand.Intn(10))
		})
		promises = append(promises, promise)
	}

	promise := New(func(resolve func(int), reject func(error)) {
		reject(fmt.Errorf("This promise was rejected"))
	})
	promises = append(promises, promise)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	values, err := MergeAll(ctx, promises...).Await(ctx)
	assert.Empty(t, values)
	assert.NotNil(t, err)
}

func TestTheFuckingBug(t *testing.T) {
	type Job struct {
		Value1 int
		Value2 int
	}

	jobs := []Job{}
	promises := []*Promise[int]{}
	for i := 0; i < 10; i++ {
		job := Job{
			Value1: rand.Intn(10),
			Value2: rand.Intn(10),
		}
		jobs = append(jobs, job)
	}

	for _, _job := range jobs {
		// For new go developers loop variables are loop scoped
		// that means that the variable _job is declared only once
		// and for each iteration the value is overrited
		// that caused a lot of bugs and will be changed in go 1.22
		// for more information see https://go.dev/blog/loopvar-preview
		job := _job
		promise := New(func(resolve func(int), reject func(error)) {
			result := job.Value1 + job.Value2
			resolve(result)
		})
		promises = append(promises, promise)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	values, err := MergeAll(ctx, promises...).Await(ctx)
	assert.NotEmpty(t, values)
	assert.Nil(t, err)
}
