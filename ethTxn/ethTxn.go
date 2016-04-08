package ethTxn

import (
	"ether_go/ethKeys"
	"ether_go/ethIpc"

	"errors"
	"fmt"
	"time"	
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// SendEthereum(sender * ethKeys.AccountKey, recipient common.Address, amountToSend int64) error
// 
// 
func SendEthereum(sender * ethKeys.AccountKey, recipient common.Address, amountToSend int64) (interface{},error) {
	var ret 	interface{}
	var zero	interface{}
	
	
	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		return zero,errors.New("error in IPC")
	}

	var TxnCount string
	err:= myEipc.Call("eth_getTransactionCount",sender.PublicKeyAsHexString(),&TxnCount)
	if err != nil {
		return zero,err	
	}
	TxnCountBytes := common.FromHex(TxnCount)
	nonce := common.ReadVarInt(TxnCountBytes)
		
	var gasPriceStr string
	err = myEipc.Call("eth_gasPrice",nil,&gasPriceStr)
	if err != nil {
		return zero,err
	}
	gasPriceBytes := common.FromHex(gasPriceStr)
	gasPrice := common.BytesToBig(gasPriceBytes)
	

	var amount	 big.Int
	amount.SetInt64(amountToSend)
	var gasLimit big.Int
	gasLimit.SetInt64(121000) // because it is a send - quite standard
	data := common.FromHex("0x")
	nt,err := types.NewTransaction(nonce, recipient, &amount,&gasLimit, gasPrice, data).SignECDSA(sender.GetKey())
	if err != nil {
		return zero,err
	}
	rlpEncodedTx, err := rlp.EncodeToBytes(nt)
	if err != nil {
	    return zero,err
	}
	rawTxn := "0x"+common.Bytes2Hex(rlpEncodedTx)
	fmt.Printf("this goes in to raw tx: %v\n", rawTxn)
	
	fmt.Println("about to send Raw")
	
	err = myEipc.Call("eth_sendRawTransaction",rawTxn,&ret)
	if err != nil {
		return zero,err
	}
	return ret,nil
}

type tx struct {
	From			string
	To				string
	Value			interface{}
	Gas				interface{}
	GasPrice	interface{}
	Data			string
}

// Estimate the gas required for a contract to run
//
func estimateGas(sender  * ethKeys.AccountKey, contract string) (big.Int,error) {
	var txStruct	tx
	
	var zero			big.Int
	
	zero.SetInt64(0)
	
	fmt.Println("Estimate Gas")
	

	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		return zero,errors.New("error in IPC")
	}

	txStruct.Data    = contract
	
	var gasLimitStr string
	err := myEipc.Call("eth_estimateGas",&txStruct,&gasLimitStr)
	if err != nil {
		return zero,err
	}
	fmt.Println("Gastimate: ",gasLimitStr)
	gasLimitBytes := common.FromHex(gasLimitStr)
	gasLimit := common.BytesToBig(gasLimitBytes)

	return *gasLimit,nil
	
}



// PostContract(sender * ethKeys.AccountKey, contract string) (interface{},error) 
// contract comes from :  common.FromHex(compiledContract.Code)
// return value is the TxHash
//
func PostContract(sender * ethKeys.AccountKey, contract []byte) (interface{},error) {
	var ret 	interface{}
	var zero	interface{}
	//
	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		return zero,errors.New("error in IPC")
	}

	var TxnCount string
	err:= myEipc.Call("eth_getTransactionCount",sender.PublicKeyAsHexString(),&TxnCount)
	if err != nil {
		return zero,err	
	}
	TxnCountBytes := common.FromHex(TxnCount)
	nonce := common.ReadVarInt(TxnCountBytes)

	var gasPriceStr string
	err = myEipc.Call("eth_gasPrice",nil,&gasPriceStr)
	if err != nil {
		return zero,err
	}
	gasPriceBytes := common.FromHex(gasPriceStr)
	gasPrice := common.BytesToBig(gasPriceBytes)

	var amount	 big.Int
	amount.SetInt64(00)
	var gasLimit big.Int
	gasLimit.SetInt64(90000000)

	newContractTx := types.NewContractCreation(nonce,&amount, &gasLimit, gasPrice, contract)
	nt,err := sender.Sign(newContractTx)
	
	rlpEncodedTx, err := rlp.EncodeToBytes(nt)
	if err != nil {
	    panic(err)
	}
	
	strTxn := "0x"+common.Bytes2Hex(rlpEncodedTx)

	gasLimit,err = estimateGas(sender,common.ToHex(contract))
	if err != nil {
		return zero,nil
	}
	
	fmt.Println("Gas Limit",gasLimit.Int64())
	
	newContractTx = types.NewContractCreation(nonce,&amount,&gasLimit , gasPrice, contract)
	nt,err = sender.Sign(newContractTx)
	
	rlpEncodedTx, err = rlp.EncodeToBytes(nt)
	if err != nil {
	    panic(err)
	}
	strTxn = "0x"+common.Bytes2Hex(rlpEncodedTx)

	fmt.Println("calling sendRaw")
	err = myEipc.Call("eth_sendRawTransaction",&strTxn, &ret)
	if err != nil {
		return zero,err
	}
	fmt.Println(ret)

	return ret,nil	
	

}

// You probably have a far better way to handle this but this 
// creates a wait loop for the transaction to enter the blockchain
//
func WaitForTxnReceipt(txn interface{}) (interface{},error) {
	var ret 	interface{}
	var zero  interface{}
	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		return zero,errors.New("error in IPC")
	}
	fmt.Println()
	count := 100
	err := errors.New("43")
	for err != nil {
		err =   myEipc.Call("eth_getTransactionReceipt", txn, &ret)
		
		if err == nil {
			break
		}
		fmt.Print(".")
		time.Sleep(500*time.Millisecond)

		if count < 0 {
			return zero,errors.New("Timeout")
			
		}
		count--
	}
	fmt.Println(ret)
	return ret,nil
}

