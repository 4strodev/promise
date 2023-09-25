package promise_test

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/4strodev/promise/pkg"
)

type Job struct {
	Value1 int
	Value2 int
}

func ExampleNew() {
	p := promise.New(func(resolve func(int), reject func(error)) {
		time.Sleep(1 * time.Second)
		resolve(1)
	})

	result, err := p.Await(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output: 1
}

func ExampleMergeAll() {
	jobs := []Job{}
	promises := []*promise.Promise[int]{}
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
		p := promise.New(func(resolve func(int), reject func(error)) {
			result := job.Value1 + job.Value2
			resolve(result)
		})
		promises = append(promises, p)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	values, err := promise.MergeAll(ctx, promises...).Await(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(values)
}

func ExampleThen() {
	p := promise.New(func(resolve func(int), reject func(error)) {
		time.Sleep(time.Millisecond * 1)
		resolve(1)
	})
	ctx := context.Background()
	newPromise := promise.Then(ctx, p, func(num int) string {
		time.Sleep(time.Millisecond * 1)
		return strconv.Itoa(num * 2)
	})
	value, err := newPromise.Await(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	// Output: 2
}
