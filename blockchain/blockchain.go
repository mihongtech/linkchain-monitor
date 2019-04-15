package main

import (
	_ "bytes"
	"flag"
	"strconv"

	"net/http"
	"os"
	_ "os/exec"
	_ "path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/log"
	_ "github.com/gethctl/utils"
	"github.com/mihongtech/linkchain/client/httpclient"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	latestBlockHeight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "latest_block_height",
		Help: "latest bitcoind block height",
	})

	blockHeight int64 = 0
	httpConfig        = &httpclient.Config{}
)

func init() {
	prometheus.MustRegister(latestBlockHeight)
}

func updateBlockNumber() {
	blockNumber, err := GetBlockNumber()
	if err != nil {
		log.Error("can't get latest block number", "err", err)
	} else {
		latestBlockHeight.Set(float64(blockNumber))
	}
}

func main() {
	var (
		listenAddr    = flag.String("addr", ":9200", "listen address")
		rpcServer     = flag.String("rpcserver", "127.0.0.1", "ethereum rpc server ip")
		rpcPort       = flag.Int("rpcport", 8082, "ethereum rpc port")
		rpcuser       = flag.String("rpcuser", "lc", "rpc user")
		rpcpassword   = flag.String("rpcpassword", "lc", "rpcpassword")
		fetchInterval = flag.Int("interval", 5, "fetch interval")
		loglevel      = flag.Int("loglevel", 3, "log level")
		metrics       = flag.String("metrics_path", "/blockchain", "the metrics path of promethus")
	)
	flag.Parse()
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(*loglevel), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	httpConfig.RPCUser = *rpcuser
	httpConfig.RPCPassword = *rpcpassword
	httpConfig.RPCServer = *rpcServer + ":" + strconv.Itoa(*rpcPort)

	go func() {
		for {
			<-time.After(time.Duration(*fetchInterval) * time.Second)
			updateBlockNumber()
		}
	}()

	http.Handle(*metrics, prometheus.Handler())
	http.ListenAndServe(*listenAddr, nil)
}
