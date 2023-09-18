package promise

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResolvedPromise(t *testing.T) {
	promise := NewPromise(func(resolve func(int), reject func(error)) {
		time.Sleep(time.Millisecond * 100)
		resolve(1)
	})
	value, err := promise.Await(context.Background())
	assert.Equal(t, 1, value)
	assert.Nil(t, err)
}

func TestUnresolvedPromise(t *testing.T) {
	// Creating a promise that doesn't resolve any value
	promise := NewPromise(func(resolve func(int), reject func(error)) {
		fmt.Println("Executing callback")
	})

	// Ensure this promise is handled anyway
	_, err := promise.Await(context.Background())
	assert.NotNil(t, err)
}

func TestPanicPromise(t *testing.T) {
	promise := NewPromise(func(resolve func(int), reject func(error)) {
		panic("A simple error")
	})

	_, err := promise.Await(context.Background())
	assert.NotNil(t, err)
}

func TestTimeoutPromise(t *testing.T) {
	promise := NewPromise(func(resolve func(int), reject func(error)) {
		time.Sleep(time.Millisecond * 100)
		panic("A simple error")
	})

	context, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()
	_, err := promise.Await(context)
	assert.NotNil(t, err)
}
