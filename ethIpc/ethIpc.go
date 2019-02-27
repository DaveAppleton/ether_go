package ethIpc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

type ethIpcHandler struct {
	ipcFileLocation string
	client          *rpc.Client
}

// This assumes that viper has been initialised in the main program
// NewEthIpc sets up the path to the Ethereum IPC file
// returns true is if exists, false otherwise
// usage:
//  myIPC = new(ethIpc)
//  if myIPC.New() != nil { ....
func NewEthIpc() (*ethIpcHandler, error) {
	eh := new(ethIpcHandler)
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	eh.ipcFileLocation = viper.GetString("IPC_PATH")
	if len(eh.ipcFileLocation) == 0 {
		eh.ipcFileLocation = usr.HomeDir + "/Library/Ethereum/geth.ipc"
	}
	_, err = os.Stat(eh.ipcFileLocation)
	if os.IsNotExist(err) {
		return nil, nil
	}
	client, err := rpc.DialIPC(context.TODO(), eh.ipcFileLocation)
	// laddr := net.UnixAddr{Net: "unix", Name: eh.ipcFileLocation}
	// conn, err := net.DialUnix("unix", nil, &laddr)
	if err != nil {
		fmt.Println("DialUnix : ", err)
		return nil, err
	}
	//defer client.Close()
	eh.client = client
	return eh, err
}
func (eh *ethIpcHandler) Close() {
	if eh.client != nil {
		eh.client.Close()
	}
}

func (eh *ethIpcHandler) EthClient() (*ethclient.Client, error) {
	if eh.client == nil {
		return nil, errors.New("Invalid RPC Client")
	}
	return ethclient.NewClient(eh.client), nil
}

// Call is a direct pass through to JSON / Client
//
// REMEMBER : args is a STRUCTURE not JSON - forget this at your peril
//
func (eh *ethIpcHandler) Call(reply interface{}, serviceMethod string, args ...interface{}) error {
	err := eh.client.Call(reply, serviceMethod, args...)
	if err != nil {
		fmt.Println("Call : ", err)
		return err
	}
	return nil
}
