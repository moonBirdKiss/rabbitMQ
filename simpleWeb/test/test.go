package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	ty := make(chan float64 , 20)
	fmt.Println(reflect.TypeOf(ty))
	v_map := make(map[string] chan float64)

	// send data
	go func() {
		fmt.Println("son1 starts")
		v_map["yes"] = make(chan float64, 10)
		for i := 0 ; i < 10; i++{
			v_map["yes"] <- float64(i)
			fmt.Println("i have send a i: ", i)
		}
		fmt.Println(v_map)
	}()

	t2 := time.NewTimer(time.Second * 1)
	<- t2.C

	// receive data
	go func() {
		fmt.Println("\t son2 starts")
		fmt.Println("\t", v_map)
		for i := 0 ; i < 11 ; i++{
			data := <-v_map["yes"]
			fmt.Println("\t i have receive a i: ", data)
		}
		fmt.Println("son2 ends")
	}()


	time1 := time.NewTimer(time.Second * 10)
	<-time1.C
	fmt.Println("kang xia le suo you")

}
