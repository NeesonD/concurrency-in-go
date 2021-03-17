package main

import (
	"fmt"
	"sync"
)

func main() {
	myPool := &sync.Pool{New: func() interface{} {
		fmt.Println("create new instance")
		return struct {
		}{}
	}}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}
