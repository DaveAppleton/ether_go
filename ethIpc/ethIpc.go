package ethIpc

import (
		"os/user"	
		"os"
		"log"
		"net"
		"net/rpc/jsonrpc"


)


type ethIpcHandler struct {
	ipcFileLocation string
}

// New sets up the path to the Ethereum IPC file
// returns true is if exists, false otherwise
// usage:
//  myIPC = new(ethIpc)
//  if myIPC.New() != nil { .... 
func NewEthIpc() *ethIpcHandler  {
	eh := new(ethIpcHandler)
  usr, err := user.Current()
	if err != nil {
	    log.Println( err )
	    return nil
	}
	eh.ipcFileLocation = usr.HomeDir+"/Library/Ethereum/geth.ipc"
	_, err = os.Stat(eh.ipcFileLocation)
	if os.IsNotExist(err) {
		return nil
	}
	return eh
}

// Call is a direct pass through to JSON / Client
//
// REMEMBER : args is a STRUCTURE not JSON - forget this at your peril
//
func  (eh *ethIpcHandler) Call(serviceMethod string, args interface{}, reply interface{}) error {
	laddr := net.UnixAddr{Net: "unix", Name: eh.ipcFileLocation}
	conn, err := net.DialUnix("unix", nil, &laddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)
	
	err = client.Call(serviceMethod,args,&reply)
	if err != nil {
		return err
	}
	return nil
}

	