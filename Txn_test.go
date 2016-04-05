package ether_go

import (
	"testing"
	"ether_go/ethKeys"
	"ether_go/ethTxn"
	"github.com/ethereum/go-ethereum/common"

)

func TestSendTxn(t * testing.T) {
	banker := ethKeys.NewKey("banker")
	err := banker.RestoreOrCreate()
	if err != nil {
		t.Error(err)
		t.Fail()	
	}
	to := common.HexToAddress("0x39c5ab6cf1fd6036505133737d4dd655fefd9d8d")
	txHash,err := ethTxn.SendEthereum(banker,to,1000000000000000000)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log("TxHash : ", txHash)
	txResult,err := ethTxn.WaitForTxnReceipt(txHash)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(txResult)
}