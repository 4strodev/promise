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
	locker    sync.Mutex
	completed bool
}

// Creates a new promise and exeuctes callback handling a possible panic
func New[T any](callback func(resolve func(T), reject func(error))) *Promise[T] {
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

func (p *Promise[T]) complete() {
	p.locker.Lock()
	p.completed = true
	p.locker.Unlock()
}

func (p *Promise[T]) isCompleted() bool {
	p.locker.Lock()
	result := p.completed
	p.locker.Unlock()
	return result
}

func (p *Promise[T]) handlePanic() {
	recovered := recover()
	if recovered != nil {
		p.reject(fmt.Errorf("Error on promise: %v", recovered))
	}
}

func (p *Promise[T]) resolve(value T) {
	if p.isCompleted() {
		return
	}
	p.value = value
	close(p.done)
	p.complete()
}

func (p *Promise[T]) reject(err error) {
	if p.isCompleted() {
		return
	}
	p.err = err
	close(p.done)
	p.complete()
}

func (p *Promise[T]) Await(ctx context.Context) (T, error) {
	go func() {
		<-ctx.Done()
		p.reject(ctx.Err())
	}()
	<-p.done
	return p.value, p.err
}
