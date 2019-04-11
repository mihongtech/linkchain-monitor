// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package scan

import (
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/linkchain/common/util/log"
	"github.com/mihongtech/linkchain-monitor/nodediscover/discover"
)

const (
	TcpTimeout = 500 * time.Millisecond
	MaxDepth   = 2
)

type NodesScan struct {
	ResultLock    sync.RWMutex
	BootNodesLock sync.RWMutex
	udp           *discover.UDP
	resultNodes   map[string]*discover.Node
	RedisAddr     string
	Bootnodes     []*discover.Node // list of bootstrap nodes
	flushDuration time.Duration
}

func NewNodesScan(udp *discover.UDP, bootnodes []*discover.Node, redisAddr string, duration time.Duration) *NodesScan {
	return &NodesScan{
		udp:           udp,
		Bootnodes:     bootnodes,
		RedisAddr:     redisAddr,
		flushDuration: duration,
	}
}

func (scan *NodesScan) findNodeCmd(node *discover.Node, target *discover.Node) ([]*discover.Node, error) {
	var result []*discover.Node = make([]*discover.Node, 0, 0)
	log.Debug("send find node to ", "node.ID", node.ID.String())
	r, e := scan.udp.Findnode(node.ID, node.ADDR(), target.ID)
	if e == nil {
		log.Debug("send find node success", "id", node.ID.String(), "addr", node.ADDR())
		for _, n := range r {
			log.Debug("find node", "id", n.ID.String(), "addr", n.ADDR())
		}
		return r, nil
	} else {
		log.Debug("find node failed", "id", node.ID.String(), "addr", node.ADDR(), "err", e)
		return r, e
	}

	return result, nil
}

func checkTCPPort(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, TcpTimeout)
	if nil != err {
		return false
	}
	defer conn.Close()
	return true
}

func (scan *NodesScan) pingCmd(node *discover.Node) error {
	log.Debug("send ping node to ", "node.ID", node.ID.String())
	err := scan.udp.Ping(node.ID, node.ADDR())
	if err == nil {
		log.Debug("send ping node success", "id", node.ID.String(), "addr", node.ADDR())
	} else {
		log.Debug("ping failed", "id", node.ID.String(), "addr", node.ADDR(), "err", err)
		return err
	}

	log.Debug("wait for pong from ", "node.ID", node.ID.String())
	scan.udp.Waitping(node.ID)

	return nil
}

func (scan *NodesScan) addNodeDataList(resultList []*discover.Node) {
	scan.ResultLock.Lock()
	defer scan.ResultLock.Unlock()
	for _, node := range resultList {
		scan.resultNodes[node.IP.String()] = node
	}
	log.Debug("result Map size", "len(scan.resultNodes)", len(scan.resultNodes))
}

func (scan *NodesScan) addNodeData(node *discover.Node) {
	scan.ResultLock.Lock()
	defer scan.ResultLock.Unlock()
	scan.resultNodes[node.IP.String()] = node
	log.Debug("result Map size", "len(scan.resultNodes)", len(scan.resultNodes))
}

func GetRandWait(base int, count int) int {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(count)
	return base + index
}

func (scan *NodesScan) runFindNodeCmd(node *discover.Node, targetNodes []*discover.Node, depth int) {
	if (depth > MaxDepth) || (len(targetNodes) == 0) {
		return
	}

	maxFailed := 5
	failed := 0
	allNeighbors := make([]*discover.Node, 0, 0)
	for _, target := range targetNodes {
		nodeList, err := scan.findNodeCmd(node, target)
		log.Debug("nodeList size", "len(nodeList)", len(nodeList))
		if (err == nil) || len(nodeList) > 0 {
			scan.addNodeDataList(nodeList)
			allNeighbors = append(allNeighbors, nodeList...)
			failed = 0
		} else {
			failed += 1
			if failed > maxFailed {
				node.ConnWait = GetRandWait(1, 5)
				return
			}
		}
	}

	// TODO: add check
	// no new neighbors found , no need to add depth
	scan.runFindNodeCmd(node, allNeighbors, depth+1)

}

func (scan *NodesScan) runFindNode(node *discover.Node) {
	if err := scan.pingCmd(node); err != nil {
		// check alive
		if checkTCPPort(node.IP.String() + ":" + string(node.TCP)) {
			node.Uptime = time.Now()
			scan.addNodeData(node)
		} else {
			log.Debug("ping cmd failed")
		}
	}

	// pingpong success add self to result
	node.Uptime = time.Now()
	scan.addNodeData(node)

	nodeList, err := scan.findNodeCmd(node, node)
	if (err == nil) || len(nodeList) > 0 {
		scan.addNodeDataList(nodeList)
		node.Uptime = time.Now()
		scan.addNodeData(node)
	} else {
		node.ConnWait = GetRandWait(1, 5)
		log.Debug("find node failed")
	}

	scan.runFindNodeCmd(node, nodeList, 1)

}

func (scan *NodesScan) addNewNodes(currNodes map[string]*discover.Node) {

	scan.ResultLock.RLock()
	defer scan.ResultLock.RUnlock()
	log.Info("currNodes", "currNodes", len(currNodes))
	log.Info("resultNodes", "resultNodes", len(scan.resultNodes))
	for k, v := range scan.resultNodes {
		currNodes[k] = v
	}
	log.Info("currNodes update", "currNodes", len(currNodes))
}

func (scan *NodesScan) addBootNodes(currNodes map[string]*discover.Node) {
	scan.BootNodesLock.RLock()
	defer scan.BootNodesLock.RUnlock()
	for _, node := range scan.Bootnodes {
		currNodes[node.IP.String()] = node
	}
}

func (scan *NodesScan) worker(nodes chan *discover.Node, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	defer wg.Done()

	for node := range nodes {
		scan.runFindNode(node)
	}
}

func (scan *NodesScan) ScanNodesLoop() {

	scan.resultNodes = make(map[string]*discover.Node)
	var currNodes = make(map[string]*discover.Node)

	startTime := time.Now()

	for {
		// add bootnodes to currnodes
		scan.addBootNodes(currNodes)
		log.Info("start to scan")
		wg := new(sync.WaitGroup)
		lCh := make(chan *discover.Node)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go scan.worker(lCh, wg)
		}

		for _, node := range currNodes {
			if node.ConnWait > 0 {
				node.ConnWait -= 1
				log.Debug("wait for node conn")
				continue
			}
			lCh <- node
		}
		close(lCh)
		log.Info("wait all scan thread over")
		wg.Wait()

		// remove timeout node
		scan.removeExpireNode(currNodes)

		scan.addNewNodes(currNodes)

		now := time.Now()
		log.Debug("check duration", "now", now, "startTime", startTime,
			"flushDuration", scan.flushDuration, "now.Sub(startTime)", now.Sub(startTime))
		if now.Sub(startTime) >= scan.flushDuration {

			scan.flushCurrData(currNodes)
			startTime = now
		}
		time.Sleep(time.Second)
		scan.ResultLock.Lock()
		scan.resultNodes = make(map[string]*discover.Node)
		scan.ResultLock.Unlock()

	}
}

func (scan *NodesScan) removeExpireNode(nodeMap map[string]*discover.Node) {
	now := time.Now()
	removeKeys := make([]string, 0, 0)
	for key, node := range nodeMap {
		if now.Sub(node.Uptime) > scan.flushDuration {
			removeKeys = append(removeKeys, key)
		}
	}

	for _, key := range removeKeys {
		delete(nodeMap, key)
	}
}
