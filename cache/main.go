package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

type Memory struct {
	f     Function
	cache map[int]FunctionResult
	lock  sync.RWMutex
}

type Function func(key int) (interface{}, error)

type FunctionResult struct {
	value interface{}
	err   error
}

// Constructor de la cache
func NewCache(f Function) *Memory {
	return &Memory{
		f:     f,
		cache: make(map[int]FunctionResult),
	}
}

// Get de la cache
func (m *Memory) Get(key int) (interface{}, error) {
	m.lock.RLock()
	result, exists := m.cache[key]
	m.lock.RUnlock()
	if !exists {
		m.lock.Lock()
		result.value, result.err = m.f(key)
		m.cache[key] = result
		m.lock.Unlock()
	}
	return result.value, result.err
}

func GetFibonacci(n int) (interface{}, error) {
	return Fibonacci(n), nil
}

func main() {
	cacheFibonacci := NewCache(GetFibonacci)
	fibo := []int{42, 40, 41, 42, 38}
	var wg sync.WaitGroup
	for _, n := range fibo {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			start := time.Now()
			value, err := cacheFibonacci.Get(index)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("%d, %s, %d\n", index, time.Since(start), value)
		}(n)
	}
	wg.Wait()
}
