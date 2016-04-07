# ether-go

This library ws created because

1. you may not like signing via RPC or IPC calls
2. you may not want to use the key store which keeps wanting passwords etc. and only workds if a key is unlocked.


You will need to load the following...

```
go get github.com/ethereum/go-ethereum/common
go get github.com/ethereum/go-ethereum/core/types
go get github.com/ethereum/go-ethereum/crypto
```

NOTE: to run some of the examples You will need to have some ether in the banker account 
(run it on the test net then you can mine into it - get the address from getKeyAddress.go)

The tests are a bit dodgy but the functionality is there....

#### Basics / Keys

Create an account if it does not exist - or load it

```
	banker := ethKeys.NewKey("banker")
	err = banker.RestoreOrCreate()
	if err != nil {
		fmt.Printf("Creating Banker %v\n",err)
		os.Exit(1)
	}
```

Get the address as a string (eg to set as a mining target or receive ether, to put in JSON)
```
	fmt.Println(banker. PublicKeyAsHexString())
```

Sign a transaction
```
	signedTxn,err := banker.Sign(unsignedTransaction)
```

#### Basics / Transactions 

PostContract - send a contract onto the blockchain

```
	hash, err := ethTxn.PostContract(banker,common.FromHex(compiledContract.Code))
```

SendEthereum - transfer ether
```
	to := common.HexToAddress("0x39c5ab6cf1fd6036505133737d4dd655fefd9d8d")
	txHash,err := ethTxn.SendEthereum(banker,to,1000000000000000000)
```

and a sample function to wait for the transaction to enter the blockchain

```
	txResult,err := ethTxn.WaitForTxnReceipt(txHash)
```




### mixed ether-go and JSON examples:


1. clientVersion.go   just JSON-RPC untidy 
2. bankerBalance.go

### ether-go examples

1. getKeyAddress.go						create or load a key and print address
2. firstContractPlayground.go	a load of calls here - have a look


 getKeyAddress.go will create a key (banker) and print out its address
 so that you can send it funds so that it can work in the next example
 

 firstContractPlayground.go contains code to 
 1. Compile a contract
 2. Post that contract to the blockchain (from the banker acount)
 
 This does, of course, require that the banker account has sufficient funds

