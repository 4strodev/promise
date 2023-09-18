package promise

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnresolvedPromise(t *testing.T) {
	// Creating a promise that doesn't resolve any value
	promise := NewPromise(func(resolve func(int), reject func(error)) {
		fmt.Println("Executing callback")
	})

	// Ensure this promise is handled anyway
	_, err := promise.Await(context.Background())
	assert.NotNil(t, err)
}

