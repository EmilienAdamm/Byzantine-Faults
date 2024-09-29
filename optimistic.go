package main

import (
	"sync/atomic"
)

type OptimisticStore struct { 
	size int 
	atomicValues []int64 
} 

func NewOptimisticStore(size int) *OptimisticStore { 
	atomicValues := make([]int64, size) 
	return &OptimisticStore{
		size: size, 
		atomicValues: atomicValues,
	} 
} 

func (s *OptimisticStore) Name() string { 
	return "Optimistic" 
} 

func (s *OptimisticStore) Size() int { 
	return s.size 
} 

func (s *OptimisticStore) At(index int) int64 { 
	return atomic.LoadInt64(&s.atomicValues[index]) 
} 

func (s *OptimisticStore) Add(index int, amount int64) {
	for {
		oldValue := s.At(index)
		newValue := oldValue + amount
		if atomic.CompareAndSwapInt64(&s.atomicValues[index], oldValue, newValue) {
			break
		}
	}
}

func (s *OptimisticStore) Substract(index int, amount int64) {
	s.Add(index, -amount)
}