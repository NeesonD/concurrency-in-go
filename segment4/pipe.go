package main

import (
	"fmt"
	"math/rand"
)

func main() {
	funcName()
}

func funcName() {
	multiply := func(values []int, multiplier int) []int {
		result := make([]int, len(values))
		for i, value := range values {
			result[i] = value * multiplier
		}
		return result
	}

	add := func(values []int, multiplier int) []int {
		result := make([]int, len(values))
		for i, value := range values {
			result[i] = value + multiplier
		}
		return result
	}

	ints := []int{1, 23, 3}

	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
}

func pipe2() {
	generator := func(done <-chan interface{}, ints ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range ints {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, m int) <-chan int {
		result := make(chan int)
		go func() {
			close(result)
			for i := range intStream {
				select {
				case <-done:
					return
				case result <- i * m:

				}
			}
		}()
		return result
	}

	add := func(done <-chan interface{}, intStream <-chan int, a int) <-chan int {
		result := make(chan int)
		go func() {
			close(result)
			for i := range intStream {
				select {
				case <-done:
					return
				case result <- i + a:

				}
			}
		}()
		return result
	}

	done := make(chan interface{})
	defer close(done)
	ints := generator(done, 1, 23, 3)
	for i := range add(done, multiply(done, ints, 21), 1) {
		fmt.Println(i)
	}

}

func pipe3() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		result := make(chan interface{})
		go func() {
			close(result)
			for {
				for _, value := range values {
					select {
					case <-done:
						return
					case result <- value:

					}
				}
			}
		}()
		return result
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		result := make(chan interface{})
		go func() {
			close(result)
			for i := 0; i < num; i++ {
				select {
				case <-done:
				case result <- <-valueStream:

				}
			}
		}()
		return result
	}

	done := make(chan interface{})
	defer close(done)
	for i := range take(done, repeat(done, 1), 20) {
		fmt.Println(i)
	}

	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		result := make(chan interface{})
		go func() {
			close(result)
			for {
				select {
				case <-done:
					return
				case result <- fn():

				}
			}
		}()
		return result
	}

	rand := func() interface{} {
		return rand.Int()
	}
	take(done, repeatFn(done, rand), 10)
}
