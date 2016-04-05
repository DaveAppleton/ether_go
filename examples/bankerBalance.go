package main

// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getbalance

import (
		"ether_go/ethKeys"
		
		"net/http"
		"strings"
		"io/ioutil"	
    "fmt"
    "encoding/json"
    "os"
    "math/big"	
)



func Call(address string, method string, id interface{}, params []interface{})(map[string]interface{},error){
    data, err := json.Marshal(map[string]interface{}{
        "method": method,
        "id":     id,
        "params": params,
    })

    if err != nil {
    	return nil, err
    }
    
    resp, err := http.Post(address, "application/json", strings.NewReader(string(data)))
    if err != nil {
    	return nil, err
    }
    
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
    	return nil, err
    }
    result := make(map[string]interface{})
    err = json.Unmarshal(body, &result)
    if err != nil {
    	return nil, err
    }

    return result, nil
}


func main(){
	banker := ethKeys.NewKey("banker")
	err := banker.RestoreOrCreate()
	if err != nil {
		fmt.Printf("Creating Banker %v\n",err)
		os.Exit(1)
	}
	fmt.Println(banker.PublicKey().Hex())
	params  := []interface{}{banker.PublicKey().Hex(),"latest"} 

	res2,err2:=Call("http://127.0.0.1:8545", "eth_getBalance", 1, params)
	if err2 != nil {
		fmt.Println("call failed - balance")
		return
	} 
	fmt.Println(res2)
	balStr := res2["result"].(string	)
	b := new(big.Int)
	b.SetString(balStr,0)
	fmt.Println("Balance ",b)


          	
}
