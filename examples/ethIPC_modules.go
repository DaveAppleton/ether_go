package main

// +build ignore

import (
	"ether_go/ethIpc"
		
	"log"
	"fmt"
)



func main() {
	var ret map [string] string
	var params interface{}

	handler := ethIpc.NewEthIpc()

	err := handler.Call("modules",&params, &ret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ret)


}