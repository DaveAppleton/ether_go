// New project target
// HelloGold creates the contract
// banker is set as the banker
// account1 and account2 get accounts in the contract
// 
// Web interface allows
//   MYR to be added to account1 & account2
//   Gold to be added to banker
//   account1 & account2 to buy gold
//   account1 & account2 to sell gold
//   balances of all accounts to be seen
//
//
package main

import (
	// GOLANG STUFF
	"log"
	"fmt"
	"io/ioutil"
	"os"

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
	contract, err := loadContractFromFile("contracts/contract.sol")
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
	
	hash, err := ethTxn.PostContract(banker,common.FromHex(compiledContract.Code))
	fmt.Println("Posted - got ",hash)
	
//	nonce := uint64( 1048600)
//	var amount	 big.Int
//	amount.SetInt64(00)
//	var gasLimit big.Int
//	gasLimit.SetInt64(90000000)
//	var gasPrice big.Int
//	gasPrice.SetInt64(20000000000)
//
//	newContractTx := types.NewContractCreation(nonce,&amount, &gasLimit, &gasPrice, common.FromHex(compiledContract.Code))
//	nt,err := banker.Sign(newContractTx)
//	
//	rlpEncodedTx, err := rlp.EncodeToBytes(nt)
//	if err != nil {
//	    panic(err)
//	}
//	
//	strTxn := "0x"+common.Bytes2Hex(rlpEncodedTx)
//	
//	fmt.Printf("this goes in to raw tx: %v\n",strTxn )
//
//
//
//	handler := ethIpc.NewEthIpc()
//	var ret interface{} // map[ string] string
//	fmt.Println("calling sendRaw")
//	err = handler.Call("eth_sendRawTransaction",&strTxn, &ret)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(ret)


}