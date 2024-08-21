package main

import "fmt"

func main() {

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for {
			ch1 <- 1
		}
	}()

	go func() {
		for {
			ch2 <- 2
		}
	}()

	for {
		select {
		case msg := <-ch1:
			fmt.Println(msg)
		case msg := <-ch2:
			fmt.Println(msg)
		default:
			fmt.Println("no message received")
		}
	}

}
