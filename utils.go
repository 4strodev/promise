package promise

import (
	"context"
	"sync"
)

func AwaitAll[T any](ctx context.Context, promises ...Promise[T]) ([]T, error) {
	mutex := sync.Mutex{}
	resolvedValues := make([]T, len(promises))
	for _, promise := range promises {
		// TODO await with context
		go func(promise Promise[T]) {
			mutex.Lock()
			value, err := promise.Await(ctx)
			if err != nil {
	
			}
			resolvedValues = append(resolvedValues, value)
		}(promise)
	}

	return resolvedValues, nil
}
