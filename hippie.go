package main

import (
	//"time"
)

type HippieStore struct { 
	size int 
	values []int64 
} 

func NewHippieStore(size int) *HippieStore { 
	return &HippieStore{ 
		size: size, values: make([]int64, size), 
	} 
} 

func (hs *HippieStore) Name() string { 
	return "Hippie" 
} 

func (hs *HippieStore) Size() int { 
	return hs.size 
} 

func (hs *HippieStore) At(index int) int64 { 
	return hs.values[index] 
}

func (hs *HippieStore) Add(index int, amount int64) { 
	currentValue := hs.values[index] 
	//time.Sleep(1)
	hs.values[index] = currentValue + amount 
} 

func (hs *HippieStore) Substract(index int, amount int64) { 
	currentValue := hs.values[index] 
	hs.values[index] = currentValue- amount 
}