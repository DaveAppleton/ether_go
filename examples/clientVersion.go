package main

// https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_clientversion

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
    "fmt"
//    "os"
)

func Call(address string, method string, id interface{}, params []interface{})(map[string]interface{},error){
    data, err := json.Marshal(map[string]interface{}{
        "method": method,
        "id":     id,
        "params": params,
    })
    fmt.Println(string(data))
    if err != nil {
        log.Fatalf("Marshal: %v", err)
    	return nil, err
    }
    resp, err := http.Post(address,
        "application/json", strings.NewReader(string(data)))
    if err != nil {
        log.Fatalf("Post: %v", err)
    	return nil, err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("ReadAll: %v", err)
    	return nil, err
    }
    result := make(map[string]interface{})
    err = json.Unmarshal(body, &result)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    	return nil, err
    }
    //log.Println(result)
    return result, nil
}

func main(){
	fmt.Println("Get Client Version")
	res, err:=Call("http://127.0.0.1:8545", "web3_clientVersion", 67, []interface{}{})
	if err!=nil{
        log.Fatalf("Err: %v", err)
	}
	log.Println("result = ",res["result"])
}