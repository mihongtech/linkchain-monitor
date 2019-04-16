package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mihongtech/linkchain-monitor/explorer/httpprotocol"
	prometheusClient "github.com/prometheus/client_golang/api/prometheus/v1"
)

func GetAuthNodeDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	client, err := getPrometheusCLient(Cfg.ip, Cfg.port)
	if err != nil {
		httpprotocol.EncodeResult(w, httpprotocol.Failed(err))
		return
	}
	nodeIP := ps.ByName("node")

	authNodeDetail, err := getAuthNodeDetail(client, nodeIP)
	if err != nil {
		httpprotocol.EncodeResult(w, httpprotocol.Failed(err))
		return
	}

	httpprotocol.ResponceSuccessWithBody(w, authNodeDetail)
	return
}

func GetNodes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	client, err := getPrometheusCLient(Cfg.ip, Cfg.port)
	if err != nil {
		httpprotocol.EncodeResult(w, httpprotocol.Failed(err))
		return
	}
	authNodes, err := getAuthNodes(client)
	if err != nil {
		httpprotocol.EncodeResult(w, httpprotocol.Failed(err))
		return
	}

	followNodes, err := getFollowNodes(client, authNodes)
	if err != nil {
		httpprotocol.EncodeResult(w, httpprotocol.Failed(err))
		return
	}

	httpprotocol.ResponceSuccessWithBody(w, map[string]interface{}{"authNodes": authNodes, "followNodes": followNodes})
	return
}

func GetLinkchainOverview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	client, err := getPrometheusCLient(Cfg.ip, Cfg.port)
	if err != nil {
		httpprotocol.EncodeResult(w, httpprotocol.Failed(err))
		return
	}

	blockNumber, err := getPrometheusBlockNumber(client)
	if err != nil {
		httpprotocol.EncodeResult(w, httpprotocol.Failed(err))
		return
	}

	// TODO: add follow node count
	httpprotocol.EncodeResult(w, httpprotocol.SucceedWithResult(
		linkchainStatus{Overview: linkchainOverview{BlockHeight: blockNumber,
			FollowNodeCount: 0},
		}))
	return
}

func getFollowNodes(client prometheusClient.API, authNodes []authNode) ([]RegionStatistics, error) {
	allNodes := GetNodeSet()
	regions := GetNodeRegions()

	followNodeCount := len(allNodes)
	for _, node := range authNodes {
		if _, ok := allNodes[node.IP]; ok {
			key := node.Region
			if "中国" != node.Country {
				key = node.Country
			}
			followNodeCount -= 1
			regions[key] = RegionStatistics{Count: regions[key].Count - 1, IPRegion: regions[key].IPRegion}
			if regions[key].Count <= 0 {
				delete(regions, key)
			}
		}
	}

	output := make([]RegionStatistics, 0, len(regions))
	for _, region := range regions {
		output = append(output, RegionStatistics{Count: region.Count, IPRegion: region.IPRegion, Percentage: float64(region.Count) / float64(followNodeCount)})
	}

	return output, nil
}
