# Go Promise package
A simple library that allows you to create promises with go.

## Installation

```sh
go get github.com/4strodev/promise
```

## Examples
Creating a promise
```go
package main

import (
	"context"
	"fmt"
	"github.com/4strodev/promise"
	"time"
)

func main() {
	p := promise.New(func(resolve func(int), reject func(error)) {
		time.Sleep(time.Millisecond * 100)
		resolve(1)
	})

	result, err := p.Await(context.Background())
	if err != nil {
		panic(err)
	}
	// Result will be 1
	fmt.Println(result)
}
```

Working with multiple promises
```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/4strodev/promise"
)

type Job struct {
	Value1 int
	Value2 int
}

func main() {
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
```

