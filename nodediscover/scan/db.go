package scan

import (
	"github.com/garyburd/redigo/redis"
	"github.com/linkchain/common/util/log"
	"github.com/mihongtech/linkchain-monitor/nodediscover/discover"
)

func GetAllDBNodes(c redis.Conn) []string {
	dbNodes := make([]string, 0, 0)

	is_key_exit, err := redis.Bool(c.Do("EXISTS", "runningNodes"))
	if err != nil {
		log.Error("check exists key failed ", "err", err)
		return dbNodes
	} else if !is_key_exit {
		return dbNodes
	}

	// get all nodes from db
	values, e := redis.Strings(c.Do("SMEMBERS", "runningNodes"))
	if e != nil {
		log.Error("get key members failed ", "err", e)
		return dbNodes
	}

	for _, value := range values {
		dbNodes = append(dbNodes, value)
	}
	return dbNodes
}

func addDBData(c redis.Conn, nodes []string) {
	for _, node := range nodes {
		_, err := c.Do("SADD", "runningNodes", node)
		if err != nil {
			log.Error("redis sadd error ", "err", err)
		}
	}
}

func delDBData(c redis.Conn, nodes []string) {
	for _, node := range nodes {
		_, err := c.Do("SREM", "runningNodes", node)
		if err != nil {
			log.Error("redis srem error ", "err", err)
		}
	}
}

func (scan *NodesScan) checkUpdate(dbNodes []string, resultMap map[string]*discover.Node) ([]string, []string) {
	addNodes := make([]string, 0, 0)
	delNodes := make([]string, 0, 0)
	dbDataMap := make(map[string]int)

	//get del data
	for _, key := range dbNodes {
		if _, ok := resultMap[key]; ok {
			// no need to check uptime
			// scanner will check the node is alive
		} else {
			delNodes = append(delNodes, key)
		}
		dbDataMap[key] = 0
	}

	//get add data
	for key, _ := range resultMap {
		if _, ok := dbDataMap[key]; !ok {
			addNodes = append(addNodes, key)
		}
	}

	return addNodes, delNodes

}

func (scan *NodesScan) flushCurrData(resultMap map[string]*discover.Node) {

	log.Info("flush data to db", "data size", len(resultMap))
	c, err := redis.Dial("tcp", scan.RedisAddr)
	if err != nil {
		log.Error("flush data to db", "data size", "err", err)
		return
	}

	dbNodes := GetAllDBNodes(c)

	// get all update
	addNodes, delNodes := scan.checkUpdate(dbNodes, resultMap)

	addDBData(c, addNodes)

	delDBData(c, delNodes)

	defer c.Close()
}
