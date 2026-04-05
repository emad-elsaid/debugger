package main

import (
	"fmt"
	"time"
)

func main() {
	var x1 string = "Hello world"
	var x2 int = 123
	var x3 float32 = 1.5
	var x4 uint = 123

	var x1s []string = []string{x1, x1, x1, x1}
	var x2s []int = []int{x2, x2, x2, x2}

	var x1a [4]string = [4]string{x1, x1, x1, x1}
	var x2a [4]int = [4]int{x2, x2, x2, x2}

	var b1 bool = true
	var b2 bool = false

	f1 := main

	p1 := &x1
	p2 := &x2
	p3 := &x3
	p4 := &x4

	m1 := map[int]string{1: "s1", 2: "s2", 3: "s3"}
	m2 := map[string]int{"s1": 1, "s2": 2, "s3": 3}

	c1 := complex(1, 2)

	ch1 := make(chan bool)
	go func() { ch1 <- true }()
	ch2 := make(chan bool, 3)
	go func() { ch2 <- true }()
	ch3 := make(chan bool)
	close(ch3)

	fmt.Println(x1, x2, x3, x4, x1s, x2s, x1a, x2a, b1, b2, f1, p1, p2, p3, p4, m1, m2, c1)
	for {
		time.Sleep(time.Second)
	}
}
