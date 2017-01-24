package ethTxn

import (
	"github.com/DaveAppleton/ether_go/ethIpc"
	"github.com/DaveAppleton/ether_go/ethKeys"

	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/net/context"
)

// SendEthereum(sender * ethKeys.AccountKey, recipient common.Address, amountToSend int64) error
//
//
func SendEthereum(sender *ethKeys.AccountKey, recipient common.Address, amountToSend int64) (interface{}, error) {
	var ret interface{}
	var zero interface{}

	myEipc, err := ethIpc.NewEthIpc()
	if err != nil {
		return zero, err
	}
	defer myEipc.Close()

	ec, _ := myEipc.EthClient()

	nonce, err := ec.NonceAt(context.TODO(), sender.PublicKey(), nil)
	gasPrice, err := ec.SuggestGasPrice(context.TODO())
	if err != nil {
		return zero, err
	}
	fmt.Println("Nonce : ", nonce)
	fmt.Println("GasPrice : ", gasPrice)
	s := types.NewEIP155Signer(params.TestnetChainConfig.ChainId)

	var amount big.Int
	amount.SetInt64(amountToSend)
	var gasLimit big.Int
	gasLimit.SetInt64(121000) // because it is a send - quite standard
	data := common.FromHex("0x")
	t := types.NewTransaction(nonce, recipient, &amount, &gasLimit, gasPrice, data)
	nt, err := types.SignTx(t, s, sender.GetKey())
	if err != nil {
		return zero, err
	}
	rlpEncodedTx, err := rlp.EncodeToBytes(nt)
	if err != nil {
		return zero, err
	}
	err = myEipc.Call(&ret, "eth_sendRawTransaction", common.ToHex(rlpEncodedTx))
	return ret, err
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

	myEipc, err := ethIpc.NewEthIpc()
	if err != nil {
		return zero, errors.New("error in IPC")
	}

	txStruct.Data = contract

	var gasLimitStr string
	err = myEipc.Call(&gasLimitStr, "eth_estimateGas", &txStruct)
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
	myEipc, err := ethIpc.NewEthIpc()
	if err != nil {
		return zero, err
	}

	ec, _ := myEipc.EthClient()

	nonce, err := ec.NonceAt(context.TODO(), sender.PublicKey(), nil)
	if err != nil {
		return zero, err
	}

	gasPrice, err := ec.SuggestGasPrice(context.TODO())
	if err != nil {
		return zero, err
	}
	var amountZero big.Int
	amountZero.SetInt64(00)
	var gasLimit big.Int
	gasLimit.SetInt64(90000000)
	cm := ethereum.CallMsg{
		From:     sender.PublicKey(),
		To:       nil,
		Gas:      &gasLimit,
		GasPrice: gasPrice,
		Value:    &amountZero,
		Data:     contract,
	}

	estGas, err := ec.EstimateGas(context.TODO(), cm)
	if err != nil {
		return zero, err
	}

	newContractTx := types.NewContractCreation(nonce, &amountZero, estGas, gasPrice, contract)
	nt, err := sender.Sign(newContractTx)

	rlpEncodedTx, err := rlp.EncodeToBytes(nt)
	if err != nil {
		panic(err)
	}

	strTxn := common.ToHex(rlpEncodedTx)

	err = myEipc.Call(&ret, "eth_sendRawTransaction", &strTxn)
	return ret, err

}

// You probably have a far better way to handle this but this
// creates a wait loop for the transaction to enter the blockchain
//
func WaitForTxnReceipt(txn interface{}) (interface{}, error) {
	var ret interface{}
	var zero interface{}
	myEipc, err := ethIpc.NewEthIpc()
	if err != nil {
		return zero, errors.New("error in IPC")
	}
	fmt.Println()
	count := 100
	err = errors.New("43")
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
