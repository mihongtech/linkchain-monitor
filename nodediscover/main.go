package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/mihongtech/linkchain-monitor/nodediscover/discover"
	"github.com/mihongtech/linkchain-monitor/nodediscover/scan"
	"github.com/mihongtech/linkchain/common/util/log"
	"github.com/mihongtech/linkchain/p2p/crypto"
)

var scanner *scan.NodesScan = nil

func startService(port int) {
	//router := NewRouter()
	//log.Info("start web service", "log", http.ListenAndServe(":"+strconv.Itoa(port), router))
}

func main() {
	var (
		listenAddr    = flag.String("addr", ":30301", "listen address")
		genKey        = flag.String("genkey", "", "generate a node key")
		nodeKeyFile   = flag.String("nodekey", "test", "private key filename")
		verbosity     = flag.Int("verbosity", int(log.LvlInfo), "log verbosity (0-9)")
		bootNodes     = flag.String("bootnodes", "enode://c95cf3128951f8bf571dc95bcc7ce70db2b23f62ab15a2f9f947d0e3c17ff25583c762117630348ba31c377a26fef269f4b7c2767a1dd6c0936c05f632411513@[122.112.249.47]:40000", "bootnodes's enode")
		redisAddress  = flag.String("redisAddress", "127.0.0.1:6379", "bootnode's enode")
		flushDuration = flag.Int("flushDuration", 3600, "input a number (seconds)")
		port          = flag.Int("port", 8080, "webservice port")

		nodeKey *ecdsa.PrivateKey
		err     error
	)
	flag.Parse()

	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(*verbosity))
	log.Root().SetHandler(glogger)

	switch {
	case *genKey != "":
		nodeKey, err = crypto.GenerateKey()
		if err != nil {
			log.Error(fmt.Sprintf("could not generate key: %v", err))
			os.Exit(1)
		}
		if err = crypto.SaveECDSA(*genKey, nodeKey); err != nil {
			log.Error(fmt.Sprintf("%v", err))
			os.Exit(1)
		}
		return
	case *nodeKeyFile == "":
		log.Error("Use -nodekey to specify a private key")
		os.Exit(1)
	case *bootNodes == "":
		log.Error("Use -bootNodes to specify a peer node")
		os.Exit(1)
	case *nodeKeyFile != "":
		if nodeKey, err = crypto.LoadECDSA(*nodeKeyFile); err != nil {
			log.Error(fmt.Sprintf("-nodekey: %v", err))
			os.Exit(1)
		}
	}

	addr, err := net.ResolveUDPAddr("udp", *listenAddr)
	if err != nil {
		log.Error(fmt.Sprintf("-ResolveUDPAddr: %v", err))
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Error(fmt.Sprintf("-ListenUDP: %v", err))
		os.Exit(1)
	}

	nodes := strings.Split(*bootNodes, ",")
	var bootNodeSlice []*discover.Node
	for _, n := range nodes {
		node, perr := discover.ParseNode(n)
		if perr != nil {
			log.Error(fmt.Sprintf("parse node error : %v", err))
		}
		log.Debug("peer node info is", "peer node.ID", node.ID.String(), "peer node.ip", node.IP)

		bootNodeSlice = append(bootNodeSlice, node)
	}

	cfg := discover.Config{
		PrivateKey: nodeKey,
	}

	udp, err := discover.ListenUDP(conn, cfg)
	if err != nil {
		log.Error(fmt.Sprintf("start udp failed: %v", err))
	}
	scanner = scan.NewNodesScan(udp, bootNodeSlice, *redisAddress, time.Duration(*flushDuration)*time.Second)

	go startService(*port)

	scanner.ScanNodesLoop()
}
