package promise

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
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

	for i := 0; i < 10; i++ {
		p := New(func(resolve func(int), reject func(error)) {
			time.Sleep(time.Millisecond * 100)
			resolve(rand.Intn(10))
		})
		promises = append(promises, p)
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

func TestThen(t *testing.T) {
	prom := New(func(resolve func(int), reject func(error)) {
		time.Sleep(time.Millisecond * 1)
		resolve(1)
	})
	ctx := context.Background()
	newProm := Then(ctx, prom, func(num int) *Promise[string] {
		time.Sleep(time.Millisecond * 1)
		return Resolve(strconv.Itoa(num))
	})
	value, err := newProm.Await(ctx)
	assert.Nil(t, err)
	assert.NotEmpty(t, value)
}

func TestRejectedThen(t *testing.T) {
	prom := New(func(resolve func(int), reject func(error)) {
		time.Sleep(time.Millisecond * 1)
		resolve(1)
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 1)
	defer cancel()

	newProm := Then(ctx, prom, func(num int) *Promise[string] {
		time.Sleep(time.Millisecond * 1)
		return Resolve(strconv.Itoa(num))
	})

	value, err := newProm.Await(ctx)
	assert.NotNil(t, err)
	assert.Empty(t, value)
}
