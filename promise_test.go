package promise

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	promise := NewPromise(func(resolve func(int), reject func(error)) {
		fmt.Println("Executing callback")
	})

	_, err := promise.Await(context.Background())
	assert.NotNil(t, err)
}

