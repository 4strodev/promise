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
	"github.com/4strodev/promise/pkg"
	"time"
)

func main() {
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
```

Working with multiple promises
```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/4strodev/promise/pkg"
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

## Suggestions are accepted
It does not mean that all suggestions will apply. It means that the suggestions will be
read and evaluated. Feel free to make any PR or Issue. Obviously, keep in mind that this is a personal project.
Which I made open source under the MIT license. I don't have much time, but I promise I will be
checking the repository 😉.

