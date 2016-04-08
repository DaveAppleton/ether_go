
// Project playground
//
package main

import (
	// GOLANG STUFF
	"log"
	"fmt"
	"io/ioutil"
	"os"
	"flag"

	// ether-go stuff
	"ether_go/ethKeys"
	"ether_go/ethTxn"


	// geth stuff
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"

)


func loadContractFromFile(fileName string) (string,error) {
	    dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "",err
	}
	return string(dat),err
}

func compileContract(source string) (*compiler.Contract,error) {
	var cc *compiler.Contract = nil
	sol, err := compiler.New("")
	if err != nil {
		fmt.Printf("solc not found: %v\n", err)
		return cc,err
	}
	fmt.Printf(" found Solididty version %v\n", sol.Version())
	
	// compile the contract
	contracts, err := sol.Compile(source)
	if err != nil {
		fmt.Printf("error compiling source. result %v: %v", contracts, err)
		return cc,err
	}
	if len(contracts) != 1 {
		fmt.Printf("one contract expected, got\n%s", len(contracts))
		return cc,err
	}
	for _, k := range contracts {
   cc = k
	}
	return cc,err
}

func main() {
	fileToCompile := flag.String("contract","contracts/contract.sol","File to compile")
	loadTxHash    := flag.String("txhash","0x","Load TxHash instead of posting")
	
	flag.Parse();
	contract, err := loadContractFromFile(*fileToCompile) //"contracts/contract.sol")
	if err != nil {
		log.Fatal(err);
	}
	compiledContract,err := compileContract(contract) 
	if err != nil {
		log.Fatal(err);
	}
	fmt.Println(	compiledContract.Code )
	
	banker := ethKeys.NewKey("banker")
	err = banker.RestoreOrCreate()
	if err != nil {
		fmt.Printf("Creating Banker %v\n",err)
		os.Exit(1)
	}
	var txr interface{}
	var hash interface{}
	if *loadTxHash == "0x" {
		hash, err = ethTxn.PostContract(banker,common.FromHex(compiledContract.Code))
		if err != nil {
			fmt.Printf("PostContract %v\n",err)
			os.Exit(1)
		}	
		fmt.Printf("Posted. If the wait for receipt times out you can repeat this command with -txhash=%v\n",hash)
		txr,err = ethTxn. WaitForTxnReceipt(hash)
	} else {
		txr,err = ethTxn. WaitForTxnReceipt(loadTxHash)
		
	}
	if err != nil {
		fmt.Printf("Wait for receipt: %v\n",err)
		os.Exit(1)
	}
	
	fmt.Println("Transaction Receipt : ",txr)
	
	ca_m,ok := txr.(map[string]interface{});
	if !ok {
		fmt.Println("cannot get map")
		os.Exit(1)
	}
	fmt.Println("Contract Address : ",ca_m["contractAddress"])
	
	
// Get the contract's Address :
	//fmt.Println("Contract Address : ",txr["contractAddress"])


}