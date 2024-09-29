package main

import (
	"sync"
)

type PessimisticStore struct { 
	size int
	values []int64
	mutex sync.Mutex
}

func NewPessimisticStore(size int) *PessimisticStore {
	return &PessimisticStore{
		size:   size,
		values: make([]int64, size),
	}
}

func (ps *PessimisticStore) Name() string {
	return "Pessimistic"
}

func (ps *PessimisticStore) Size() int {
	return ps.size
}

func (ps *PessimisticStore) At(index int) int64 {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	return ps.values[index]
}

func (ps *PessimisticStore) Add(index int, amount int64) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	ps.values[index] += amount
}

func (ps *PessimisticStore) Substract(index int, amount int64) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	ps.values[index] -= amount
}