package main

import (
	"fmt"
	"ether_go/ethKeys"
	"os"
)


func main() {
	banker := ethKeys.NewKey("banker")
	err := banker.RestoreOrCreate()
	if err != nil {
		fmt.Printf("Creating Banker %v\n",err)
		os.Exit(1)
	}
	fmt.Printf("banker pub key : %x\n",banker.PublicKey())
}