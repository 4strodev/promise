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
		promise := NewPromise(func(resolve func(int), reject func(error)) {
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
		promise := NewPromise(func(resolve func(int), reject func(error)) {
			time.Sleep(time.Millisecond * 100)
			resolve(rand.Intn(10))
		})
		promises = append(promises, promise)
	}

	promise := NewPromise(func(resolve func(int), reject func(error)) {
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

	// Don't use range it behaves like if it uses a pointer saves a pointer and cause many bugs
	// If you prefer to user range save job values before create promise
	/*
	for _, job := range jobs {
		val1, val2 := job.Value1, job.Value2
		promise := NewPromise(func(resolve func(int), reject func(error)) {
			result := val1 + val2
			resolve(result)
		})
		promises = append(promises, promise)
	}
	*/
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		promise := NewPromise(func(resolve func(int), reject func(error)) {
			result := job.Value1 + job.Value2
			resolve(result)
		})
		promises = append(promises, promise)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	values, err := MergeAll(ctx, promises...).Await(ctx)
	fmt.Println(values)
	assert.NotEmpty(t, values)
	assert.Nil(t, err)
}
