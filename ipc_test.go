package ether_go

import (
	"testing"

	"github.com/DaveAppleton/ether_go/ethIpc"
)

func TestNew(t *testing.T) {
	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		t.Error("Either geth is not running or we cannot initialise")
		t.Fail()
	}
}

func TestCall(t *testing.T) {
	var reply string
	var params []interface{}

	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		t.Error("cannot create")
		t.Fail()
	}
	if myEipc.Call("net_version", params, &reply) != nil {
		t.Error("Cannot get net version")
		t.Fail()
	}
	if reply != "2" {
		t.Error("net version != 2")
		t.Fail()
	}
	t.Log("reply : ", reply)
}

func TestGetTx(t *testing.T) {
	var reply interface{}
	params := "0x00fee8db65cf6d45bd874d1184e9036b1bc178c7d509dde40a23ad7bafd64b20"
	myEipc := ethIpc.NewEthIpc()
	if myEipc == nil {
		t.Error("cannot create")
		t.Fail()
	}
	if myEipc.Call("eth_getTransactionReceipt", params, &reply) != nil {
		t.Error("Cannot get Tx Receipt")
		t.Fail()
	}
	t.Log(reply)

}
