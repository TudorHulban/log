package log

import (
	"errors"
	"sync"
	"time"
)

type item[T any] struct {
	instance T
	isInUse  bool
}

type Pool[T any] struct {
	pool map[int64]*item[T]
	mu   sync.Mutex
}

func NewPool[T any]() *Pool[T] {
	return &Pool[T]{
		pool: make(map[int64]*item[T]),
	}
}

func (p *Pool[T]) Get() *T {
	p.mu.Lock()
	defer p.mu.Unlock()

	for key, value := range p.pool {
		if value.isInUse {
			continue
		}

		p.pool[key] = &item[T]{
			instance: value.instance,
			isInUse:  true,
		}

		return &value.instance
	}

	var t T

	p.pool[time.Now().UnixNano()] = &item[T]{
		instance: t,
		isInUse:  true,
	}

	return &t
}

func (p *Pool[T]) Put(back *T) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for key, value := range p.pool {
		if !value.isInUse {
			continue
		}

		p.pool[key] = &item[T]{
			instance: *back,
			isInUse:  false,
		}

		return nil
	}

	return errors.New(
		"could not return value back",
	)
}
