package ethTxn

import (
	"github.com/DaveAppleton/ether_go/ethIpc"
	"github.com/DaveAppleton/ether_go/ethKeys"

	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

// SendEthereum(sender * ethKeys.AccountKey, recipient common.Address, amountToSend int64) error
//
//
func SendEthereum(sender *ethKeys.AccountKey, recipient common.Address, amountToSend int64) (interface{}, error) {
	var ret interface{}
	var zero interface{}

	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		return zero, errors.New("error in IPC")
	}

	var TxnCount string

	err := myEipc.Call(&TxnCount, "eth_getTransactionCount", sender.PublicKeyAsHexString(), "latest")
	if err != nil {
		fmt.Println("eth_getTransactionCount ", err)
		return zero, err
	}
	TxnCountBytes := common.FromHex(TxnCount)
	nonce := common.ReadVarInt(TxnCountBytes)
	fmt.Println("Nonce : ", nonce)
	var gasPriceStr string
	err = myEipc.Call(&gasPriceStr, "eth_gasPrice", nil)
	if err != nil {
		return zero, err
	}
	gasPriceBytes := common.FromHex(gasPriceStr)
	gasPrice := common.BytesToBig(gasPriceBytes)

	var amount big.Int
	amount.SetInt64(amountToSend)
	var gasLimit big.Int
	gasLimit.SetInt64(121000) // because it is a send - quite standard
	data := common.FromHex("0x")
	t := types.NewTransaction(nonce, recipient, &amount, &gasLimit, gasPrice, data)

	s := types.NewEIP155Signer(params.TestnetChainConfig.ChainId)
	nt, err := types.SignTx(t, s, sender.GetKey())

	if err != nil {
		return zero, err
	}
	rlpEncodedTx, err := rlp.EncodeToBytes(nt)
	if err != nil {
		return zero, err
	}
	rawTxn := common.ToHex(rlpEncodedTx)
	fmt.Printf("this goes in to raw tx: %v\n", rawTxn)

	fmt.Println("about to send Raw")

	err = myEipc.Call(&ret, "eth_sendRawTransaction", rawTxn)
	if err != nil {
		fmt.Println("Sending Raw Txn ", err)
		return zero, err
	}
	return ret, nil
}

type tx struct {
	From     string
	To       string
	Value    interface{}
	Gas      interface{}
	GasPrice interface{}
	Data     string
}

// Estimate the gas required for a contract to run
//
func estimateGas(sender *ethKeys.AccountKey, contract string) (big.Int, error) {
	var txStruct tx

	var zero big.Int

	zero.SetInt64(0)

	fmt.Println("Estimate Gas")

	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		return zero, errors.New("error in IPC")
	}

	txStruct.Data = contract

	var gasLimitStr string
	err := myEipc.Call(&gasLimitStr, "eth_estimateGas", &txStruct)
	if err != nil {
		return zero, err
	}
	fmt.Println("Gastimate: ", gasLimitStr)
	gasLimitBytes := common.FromHex(gasLimitStr)
	gasLimit := common.BytesToBig(gasLimitBytes)

	return *gasLimit, nil

}

// PostContract - stcik a contract onto the blockchain
//
// would seriously suggest using Abigen to create Dapp bindings instead of uaing this
//
func PostContract(sender *ethKeys.AccountKey, contract []byte) (interface{}, error) {
	var ret interface{}
	var zero interface{}
	//
	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		return zero, errors.New("error in IPC")
	}

	var TxnCount string
	err := myEipc.Call("eth_getTransactionCount", sender.PublicKeyAsHexString(), &TxnCount)
	if err != nil {
		return zero, err
	}
	TxnCountBytes := common.FromHex(TxnCount)
	nonce := common.ReadVarInt(TxnCountBytes)

	var gasPriceStr string
	err = myEipc.Call(&gasPriceStr, "eth_gasPrice", nil)
	if err != nil {
		return zero, err
	}
	gasPriceBytes := common.FromHex(gasPriceStr)
	gasPrice := common.BytesToBig(gasPriceBytes)

	var amount big.Int
	amount.SetInt64(00)
	var gasLimit big.Int
	gasLimit.SetInt64(90000000)

	newContractTx := types.NewContractCreation(nonce, &amount, &gasLimit, gasPrice, contract)
	nt, err := sender.Sign(newContractTx)

	rlpEncodedTx, err := rlp.EncodeToBytes(nt)
	if err != nil {
		panic(err)
	}

	strTxn := common.ToHex(rlpEncodedTx)

	gasLimit, err = estimateGas(sender, common.ToHex(contract))
	if err != nil {
		return zero, nil
	}

	fmt.Println("Gas Limit", gasLimit.Int64())

	newContractTx = types.NewContractCreation(nonce, &amount, &gasLimit, gasPrice, contract)
	nt, err = sender.Sign(newContractTx)

	rlpEncodedTx, err = rlp.EncodeToBytes(nt)
	if err != nil {
		panic(err)
	}

	strTxn = common.ToHex(rlpEncodedTx)

	fmt.Println("calling sendRaw")
	err = myEipc.Call(&ret, "eth_sendRawTransaction", &strTxn)
	if err != nil {
		return zero, err
	}
	fmt.Println(ret)

	return ret, nil

}

// You probably have a far better way to handle this but this
// creates a wait loop for the transaction to enter the blockchain
//
func WaitForTxnReceipt(txn interface{}) (interface{}, error) {
	var ret interface{}
	var zero interface{}
	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		return zero, errors.New("error in IPC")
	}
	fmt.Println()
	count := 100
	err := errors.New("43")
	for err != nil {
		err = myEipc.Call(&ret, "eth_getTransactionReceipt", txn)

		if err == nil {
			break
		}
		fmt.Print(".")
		time.Sleep(500 * time.Millisecond)

		if count < 0 {
			return zero, errors.New("Timeout")

		}
		count--
	}
	fmt.Println(ret)
	return ret, nil
}
