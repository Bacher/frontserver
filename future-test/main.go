package main

import (
	"frontserver/future"
	"log"
	"time"
)

func main() {
	fut := future.New()

	for i := 0; i < 10; i++ {
		ii := i
		go func() {
			fut.Then()
			log.Println("Success", ii)
		}()
	}

	fut.Done()

	time.Sleep(2 * time.Second)
}
