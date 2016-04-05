Ether-go

You will need to load the following...

go get github.com/ethereum/go-ethereum/common
go get github.com/ethereum/go-ethereum/core/types
go get github.com/ethereum/go-ethereum/crypto

NOTE: to run some of the examples You will need to have some ether in the banker account 
(run it on the test net then you can mine into it - get the address from getKeyAddress.go)

The tests are a bit dodgy but the functionality is there....

mixed ether-go and JSON examples:
=================================

bankerBalance.go   a mix of ether-go and JSON-RPC (which has not been tidied up)
clientVersion.go   just JSON-RPC untidy 

ether-go examples
=================

getKeyAddress.go						create or load a key and print address
firstContractPlayground.go	a load of calls here - have a look


