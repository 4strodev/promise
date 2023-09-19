package promise

import (
	"context"
)

// Merge all promises and return a promise that will return an array of value
// those values will be the values obtained from promises
func MergeAll[T any](ctx context.Context, promises ...*Promise[T]) *Promise[[]T] {

	resultPromise := NewPromise(func(resolve func([]T), reject func(error)) {
		resolvedValuesChannel := make(chan T, len(promises))
		resolvedValues := make([]T, 0, len(promises))

		for _, promise := range promises {
			go func(promise *Promise[T]) {
				value, err := promise.Await(ctx)
				if err != nil {
					reject(err)
				}
				resolvedValuesChannel <- value
			}(promise)
		}

		for value := range resolvedValuesChannel {
			resolvedValues = append(resolvedValues, value)
			if len(resolvedValues) == len(promises) {
				close(resolvedValuesChannel)
			}
		}

		resolve(resolvedValues)

	})

	return resultPromise
}
