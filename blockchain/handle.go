package main

import (
	"encoding/json"
	"github.com/mihongtech/linkchain/client/httpclient"
	"github.com/mihongtech/linkchain/common/util/log"
	"github.com/mihongtech/linkchain/rpc/rpcjson"
	"github.com/mihongtech/linkchain/rpc/rpcobject"
)

func GetBlockNumber() (uint64, error) {
	method := "getBlockChainInfo"

	//call
	out, err := rpc(method, nil)
	if err != nil {
		return 0, err
	}
	info := rpcobject.ChainRSP{}
	if err := json.Unmarshal(out, &info); err != nil {
		return 0, err
	}
	return uint64(info.Chains.BestHeight), err
}

//rpc call
func rpc(method string, cmd interface{}) ([]byte, error) {
	//param
	s, _ := rpcjson.MarshalCmd(1, method, cmd)
	//log.Info(method, "req", string(s))

	//response
	rawRet, err := httpclient.SendPostRequest(s, httpConfig)
	if err != nil {
		log.Error(method, "error", err)
		return nil, err
	}

	//log.Info(method, "rsp", string(rawRet))

	return rawRet, nil
}
