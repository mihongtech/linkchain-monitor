package main

import (
	_ "context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/mihongtech/linkchain-monitor/nodediscover/scan"
	"github.com/mihongtech/linkchain/common/util/log"
)

type Statistics struct {
	RpcPort  int
	Internal int
}

type IPResult struct {
	Data IPRegion `json:"data"`
}

type IPRegion struct {
	Country   string  `json:"country"`
	Region    string  `json:"regionName"`
	City      string  `json:"city"`
	Status    string  `json:"status"`
	Longitude float32 `json:"lon"`
	Latitude  float32 `json:"lat"`
}

type RegionStatistics struct {
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
	IPRegion
}

var (
	blockHeight int64 = 0

	lock        sync.RWMutex
	Exists                                  = struct{}{}
	regionCache map[string]IPRegion         = make(map[string]IPRegion)
	regions     map[string]RegionStatistics = make(map[string]RegionStatistics)
	nodeSet     map[string]struct{}         = make(map[string]struct{})
	regionLock  sync.RWMutex
)

func GetNodeRegion(node string) (IPRegion, error) {
	regionLock.RLock()
	region, ok := regionCache[node]
	regionLock.RUnlock()
	if !ok {
		var err error
		region, err = getRegion(node)
		if err != nil {
			log.Error("get ip region failed", "err", err, "node", node)
			return IPRegion{}, err
		}
		regionLock.Lock()
		regionCache[node] = region
		regionLock.Unlock()
	}

	return region, nil
}

func GetNodeRegions() map[string]RegionStatistics {
	lock.RLock()
	defer lock.RUnlock()
	return regions
}

func GetNodeSet() map[string]struct{} {
	lock.RLock()
	defer lock.RUnlock()
	return nodeSet
}

func getRegion(node string) (IPRegion, error) {
	resp, err := http.Get("http://ip-api.com/json/" + node + "?lang=zh-CN")
	if err != nil {
		return IPRegion{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return IPRegion{}, err
	}

	var result IPRegion
	err = json.Unmarshal(body, &result)
	if err != nil {
		return IPRegion{}, errors.New("Get ip check result error")
	}

	return result, nil
}

func updateNodeCount(redisClient redis.Conn) {
	regionLock.Lock()
	defer regionLock.Unlock()
	nodes := scan.GetAllDBNodes(redisClient)
	regions = make(map[string]RegionStatistics)
	nodeSet = make(map[string]struct{})
	for _, node := range nodes {
		nodeSet[node] = Exists
		region, ok := regionCache[node]
		if !ok {
			var err error
			region, err = getRegion(node)
			if err != nil {
				log.Error("get ip region failed", "err", err, "node", node)
				continue
			}
			regionCache[node] = region
		}
		log.Debug("region is", "region", region, "node", node)
		key := region.Region
		if "中国" != region.Country {
			key = region.Country
		}
		find, ok := regions[key]
		if ok {
			regions[key] = RegionStatistics{Count: find.Count + 1, IPRegion: region}
		} else {
			regions[key] = RegionStatistics{Count: 1, IPRegion: region}
		}
		log.Debug("regions is", "regions", regions)
	}
}

func (s *Statistics) Serve() {
	redisClient, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Error("get redis cleint failed", "err", err)
		return
	}
	updateNodeCount(redisClient)

	go func() {
		for {
			<-time.After(time.Hour)
			updateNodeCount(redisClient)
		}
	}()
}
