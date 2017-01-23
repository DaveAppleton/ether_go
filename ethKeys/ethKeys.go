package ethKeys

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

type AccountKey struct {
	Name string
	Key  *ecdsa.PrivateKey
}

// NewKey creates a new named key
// simple use of named keys is to save to files
// more work required
func NewKey(name string) *AccountKey {
	ak := new(AccountKey)
	ak.Name = name
	return ak
}

// SaveKey - save key as file in current working directory
func (ac *AccountKey) SaveKey() error {
	err := crypto.SaveECDSA(ac.Name, ac.Key)
	return err
}

// LoadKey - Load key (if it exists) from current directory
func (ac *AccountKey) LoadKey() error {
	key, err := crypto.LoadECDSA(ac.Name)
	if err == nil {
		ac.Key = key
	}
	return err
}

// GenerateKey - generate a new key pair
// WARNING : will overwrite existing key pair
func (ac *AccountKey) GenerateKey() error {
	key, err := crypto.GenerateKey()
	if err != nil {
		return err
	}
	ac.Key = key
	return nil
}

// RestoreOrCreate - if exists, load key
//                 - otherwise generate new key
func (ac *AccountKey) RestoreOrCreate() error {
	err := ac.LoadKey()
	if err == nil {
		return nil
	}
	err = ac.GenerateKey()
	if err != nil {
		return err
	}
	return ac.SaveKey()
}

// GetKey - get key, e.g. for signing
func (ac *AccountKey) GetKey() *ecdsa.PrivateKey {
	return ac.Key
}

// PublicKey - Get Public Key as an address
func (ac *AccountKey) PublicKey() common.Address {

	return crypto.PubkeyToAddress(ac.Key.PublicKey)
}

// PublicKeyAsHexString - Get Public Key as a string - e.g. "0xabcd....."
func (ac *AccountKey) PublicKeyAsHexString() string {
	return crypto.PubkeyToAddress(ac.Key.PublicKey).Hex()
}

// Sign - a transaction with the key
func (ac *AccountKey) Sign(t *types.Transaction) (tr *types.Transaction, err error) {
	// see accounts/account_manager.Sign
	// crypto.Sign(hash,key)
	//params.MainnetChainConfig
	s := types.NewEIP155Signer(params.TestnetChainConfig.ChainId)
	tr, err = types.SignTx(t, s, ac.Key)

	// tr, err = t.SignECDSA(types.HomesteadSigner{}, ac.Key)
	return
}
