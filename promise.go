package promise

import (
	"context"
	"fmt"
	"sync"
)

type Promise[T any] struct {
	value     T
	err       error
	done      chan struct{}
	completed bool
}

func NewPromise[T any](callback func(resolve func(T), reject func(error))) *Promise[T] {
	promise := new(Promise[T])
	promise.completed = false
	promise.done = make(chan struct{})
	waitGroup := sync.WaitGroup{}

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		defer promise.handlePanic()
		callback(promise.resolve, promise.reject)
	}()

	go func() {
		waitGroup.Wait()
		promise.reject(fmt.Errorf("Callback finished but promise not completed"))
	}()

	return promise
}

func (self *Promise[T]) handlePanic() {
	recovered := recover()
	if recovered != nil {
		self.reject(fmt.Errorf("Error on promise: %v", recovered))
	}
}

func (self *Promise[T]) resolve(value T) {
	if self.completed {
		return
	}
	self.value = value
	close(self.done)
	self.completed = true
}

func (self *Promise[T]) reject(err error) {
	if self.completed {
		return
	}
	self.err = err
	close(self.done)
	self.completed = true
}

func (self *Promise[T]) Await(ctx context.Context) (T, error) {
	go func() {
		<-ctx.Done()
		self.reject(fmt.Errorf("Context finished"))
	}()
	<-self.done
	return self.value, self.err
}
