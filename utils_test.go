package promise

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAwaitAllResolved(t *testing.T) {
	promises := make([]*Promise[int], 0, 10)
	for i := 0; i < 10; i++ {
		promise := NewPromise(func(resolve func(int), reject func(error)) {
			time.Sleep(time.Millisecond * 100)
			resolve(10)
		})
		promises = append(promises, promise)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()
	values, err := AwaitAll(ctx, promises...).Await(ctx)
	assert.NotEmpty(t, values)
	assert.Nil(t, err)
}
