package ether_go

import (
	"testing"
	
	"ether_go/ethKeys"
	
	"math/big"
	
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	
	
)

func TestNewKey(t *testing.T) {
	myKeys :=ethKeys. NewKey("mikey")
	if myKeys == nil {
		t.Error("no keys created")
		t.Fail()
	}
	if myKeys.Name != "mikey" {
		t.Error("name not set")
		t.Fail()
	}
	if myKeys.Key != nil {
		t.Error("Key has a value")
		t.Fail()
	}
}

// Assume that banker file has been loaded in current directory & bonker has not
func TestLoadKey(t *testing.T) {
	bpk := "0xe7301a7b34e6607020e1515f381ab4d8dd484bfd"
	bKeys := ethKeys.NewKey("banker")
	err := bKeys.LoadKey()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if bKeys.PublicKey() != common.HexToAddress(bpk) {
		t.Error("Public Key does not match")
		t.Error( bKeys.PublicKey())
		t.Error( common.HexToAddress(bpk) )
		t.Fail()
	}
	mKeys := ethKeys.NewKey("Non existant file")
	err = mKeys.LoadKey()
	if err == nil {
		t.Error("did not return error for file not found")
		t.Fail()
	}
}

func TestSign(t * testing.T) {
	bKeys := ethKeys.NewKey("banker")
	err := bKeys.LoadKey()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	nonce := uint64( 1048582)
	var amount	 big.Int
	amount.SetInt64(1000000000000000000)
	var gasLimit big.Int
	gasLimit.SetInt64(90000)
	var gasPrice big.Int
	gasPrice.SetInt64(10000000000000)
	toAddress := common.HexToAddress("0x39c5ab6cf1fd6036505133737d4dd655fefd9d8d")
	data := common.FromHex("0x")
	nt,err := types.NewTransaction(nonce, toAddress, &amount,&gasLimit, &gasPrice, data).SignECDSA(bKeys.GetKey())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	// Note - 1.4 codebase ValidateSignatureValues requires extra param : homestead 
	v,r,s := nt.SignatureValues()
	if !crypto.ValidateSignatureValues(v, r, s,true) {
		t.Error("Validate Singapore Fails")
		t.Fail()
	}

}