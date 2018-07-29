package main

import (
	"math/rand"
	"time"
	"fmt"
)

//generate random number

func main() {
	
	mynum := MyRandom(10000)

	fmt.Println("Rand for 10000 is: ",mynum)
}

func MyRandom(num int) int {
	source := rand.NewSource(time.Now().Unix())
	r := rand.New(source)

	return r.Intn(num)
}
