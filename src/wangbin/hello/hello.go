package main

import (
	"container/list"
	"fmt"
	"wangbin/mymath"
)

func main() {

	// var a int
	var str string
	var l list.List

	l.PushBack(4)

	str = "goog job"
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	fmt.Println(str)
	fmt.Print("hello world!1")
	fmt.Printf("Hello, world.  Sqrt(2) = %v\n", mymath.Sqrt(2))
	fmt.Printf("Hello, world.  Sqrt1(2) = %v\n", mymath.Sqrt1(2))
	fmt.Printf("new line")
}
