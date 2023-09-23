package main

import (
	"context"
	"fmt"
	"github.com/4strodev/promise"
	"time"
)

func main() {
	promise := promise.NewPromise(func(resolve func(int), reject func(error)) {
		time.Sleep(time.Millisecond * 100)
		resolve(1)
	})

	result, err := promise.Await(context.Background())
	if err != nil {
		panic(err)
	}
	// Result will be 1
	fmt.Println(result)
}
