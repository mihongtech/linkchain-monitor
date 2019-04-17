package main

import (
	"context"
	"errors"
	"github.com/mihongtech/linkchain/common/util/log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	prometheusClient "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type config struct {
	ip      string
	rpcPort int
	port    int
	timeout time.Duration
}

var Cfg *config = nil

type linkchainStatus struct {
	Overview linkchainOverview `json:"overview"`
}

type linkchainOverview struct {
	BlockHeight     *big.Int `json:"blockHeight"`
	AuthNodeCount   int      `json:"authNodeCount"`
	FollowNodeCount int      `json:"followNodeCount"`
}

type authNodeDetail struct {
	BaseInfo      authNodeBaseInfo  `json:"baseInfo"`
	LinkchainInfo authNodeChainInfo `json:"linkchainInfo"`
}

type authNodeBaseInfo struct {
	IP          string  `json:"ip"`
	OS          string  `json:"os"`
	SysTime     int     `json:"sysTime"`
	RunningTime int     `json:"runningTime"`
	GethVersion string  `json:"gethVersion"`
	DiskUsage   float64 `json:"diskUsage"`
	MemUsage    float64 `json:"memUsage"`
	CPUUsage    float64 `json:"cpuUsage"`
}

type authNodeChainInfo struct {
	BlockHeight int64 `json:"blockHeight"`
	BlockDiff   int64 `json:"blockDiff"`
}

type authNode struct {
	IP          string `json:"ip"`
	BlockHeight int64  `json:"blockHeight"`
	IPRegion
}

func getPrometheusCLient(ip string, port int) (prometheusClient.API, error) {
	client, err := api.NewClient(api.Config{Address: "http://" + ip + ":" + strconv.Itoa(port)})
	if err != nil {
		return nil, err
	}

	return prometheusClient.NewAPI(client), nil
}

func getPrometheusBlockNumber(client prometheusClient.API) (*big.Int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Cfg.timeout*time.Second)
	defer cancel()

	start := time.Now()
	value, err := client.Query(ctx, "latest_block_height", start)
	if err != nil {
		return nil, 0, err
	}
	end := time.Now()
	if end.Sub(start) > Cfg.timeout*time.Second {
		return nil, 0, errors.New("Get Prometheus block number timeout")
	}

	if model.ValVector != value.Type() {
		return nil, 0, errors.New("Get Prometheus block number data type error")
	}

	if len(value.(model.Vector)) == 0 {
		return nil, 0, errors.New("Get Prometheus block number error")
	}

	blockNumber := big.NewInt(0)
	allNodes := GetNodeSet()
	followNodeCount := len(allNodes)
	for _, data := range value.(model.Vector) {
		number := int64(data.Value)
		if blockNumber.Cmp(big.NewInt(number)) <= 0 {
			blockNumber = big.NewInt(number)
		}
		instance := string(data.Metric["instance"])
		if !strings.Contains(instance, ":") {
			continue
		}
		ip := strings.Split(instance, ":")[0]
		if _, ok := allNodes[ip]; ok {
			followNodeCount -= 1
		}
	}
	return blockNumber, followNodeCount, nil
}

func getOSVersion(ip string) (string, error) {
	//nc := &nodeclient.Client{Host: "http://" + ip + ":" + strconv.Itoa(Cfg.ctlPort)}
	return "Have not get OS Version", nil
}

func getGethVersion(ip string) (string, error) {
	//nc := &nodeclient.Client{Host: "http://" + ip + ":" + strconv.Itoa(Cfg.ctlPort)}
	return "Have not get Geth Version", nil
}

func getAuthNodes(client prometheusClient.API) ([]authNode, error) {
	var output []authNode = make([]authNode, 0, 0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	start := time.Now()
	value, err := client.Query(ctx, "latest_block_height", start)
	if err != nil {
		return output, err
	}
	end := time.Now()
	if end.Sub(start) > time.Second {
		return output, errors.New("Get Prometheus block number timeout")
	}

	if model.ValVector != value.Type() {
		return output, errors.New("Get Prometheus block number data type error")
	}

	if len(value.(model.Vector)) == 0 {
		return output, errors.New("Get Prometheus block number error")
	}
	for _, data := range value.(model.Vector) {
		number := int64(data.Value)

		instance := string(data.Metric["instance"])
		if !strings.Contains(instance, ":") {
			continue
		}
		ip := strings.Split(instance, ":")[0]
		region, err := GetNodeRegion(ip)
		if err != nil {
			log.Error("get region failed", "err", err)
			continue
		}

		output = append(output, authNode{ip, number, region})
	}

	if len(output) == 0 {
		return output, errors.New("Can not find AuthNode")
	}
	return output, nil
}

func getNodeSysTime(client prometheusClient.API, nodeIP string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	start := time.Now()
	sysTime := 0
	value, err := client.Query(ctx, "sys_time{instance=\""+nodeIP+":80\"}", start)
	if err != nil {
		return sysTime, err
	}
	end := time.Now()
	if end.Sub(start) > time.Second {
		return sysTime, errors.New("Get Prometheus sysTime timeout")
	}

	if model.ValVector != value.Type() {
		return sysTime, errors.New("Get Prometheus sysTime data type error")
	}

	if len(value.(model.Vector)) == 0 {
		return sysTime, errors.New("Get Prometheus sysTime error")
	}

	for _, data := range value.(model.Vector) {
		sysTime = int(data.Value)
	}

	return sysTime, nil
}

func getNodeMemUsage(client prometheusClient.API, nodeIP string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	start := time.Now()
	memUsage := float64(0)
	value, err := client.Query(ctx, "mem_usage{instance=\""+nodeIP+":80\"}", start)
	if err != nil {
		return memUsage, err
	}
	end := time.Now()
	if end.Sub(start) > time.Second {
		return memUsage, errors.New("Get Prometheus memUsage timeout")
	}

	if model.ValVector != value.Type() {
		return memUsage, errors.New("Get Prometheus memUsage data type error")
	}

	if len(value.(model.Vector)) == 0 {
		return memUsage, errors.New("Get Prometheus memUsage  error")
	}

	for _, data := range value.(model.Vector) {
		memUsage = float64(data.Value)
	}

	return memUsage, nil
}

func getNodeCpuUsage(client prometheusClient.API, nodeIP string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	start := time.Now()
	cpuUsage := float64(0)
	value, err := client.Query(ctx, "cpu_usage{instance=\""+nodeIP+":80\"}", start)
	if err != nil {
		return cpuUsage, err
	}
	end := time.Now()
	if end.Sub(start) > time.Second {
		return cpuUsage, errors.New("Get Prometheus cpuUsage timeout")
	}

	if model.ValVector != value.Type() {
		return cpuUsage, errors.New("Get Prometheus cpuUsage data type error")
	}

	if len(value.(model.Vector)) == 0 {
		return cpuUsage, errors.New("Get Prometheus cpuUsage  error")
	}

	for _, data := range value.(model.Vector) {
		cpuUsage = float64(data.Value)
	}

	return cpuUsage, nil
}

func getNodeDiskUsage(client prometheusClient.API, nodeIP string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	start := time.Now()
	diskUsage := float64(0)
	value, err := client.Query(ctx, "disk_usage{instance=\""+nodeIP+":80\"}", start)
	if err != nil {
		return diskUsage, err
	}
	end := time.Now()
	if end.Sub(start) > time.Second {
		return diskUsage, errors.New("Get Prometheus diskUsage timeout")
	}

	if model.ValVector != value.Type() {
		return diskUsage, errors.New("Get Prometheus diskUsage data type error")
	}

	if len(value.(model.Vector)) == 0 {
		return diskUsage, errors.New("Get Prometheus diskUsage  error")
	}

	for _, data := range value.(model.Vector) {
		diskUsage = float64(data.Value)
	}

	return diskUsage, nil
}

func getNodeUPTime(client prometheusClient.API, nodeIP string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	start := time.Now()
	upTime := 0
	value, err := client.Query(ctx, "up_time{instance=\""+nodeIP+":80\"}", start)
	if err != nil {
		return upTime, err
	}
	end := time.Now()
	if end.Sub(start) > time.Second {
		return upTime, errors.New("Get Prometheus upTime timeout")
	}

	if model.ValVector != value.Type() {
		return upTime, errors.New("Get Prometheus upTime data type error")
	}

	if len(value.(model.Vector)) == 0 {
		return upTime, errors.New("Get Prometheus upTime error")
	}

	for _, data := range value.(model.Vector) {
		upTime = int(data.Value)
	}

	return upTime, nil
}

func getAuthNodeDetail(client prometheusClient.API, nodeIP string) (*authNodeDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	start := time.Now()
	value, err := client.Query(ctx, "latest_block_height", start)
	if err != nil {
		return nil, err
	}
	end := time.Now()
	if end.Sub(start) > time.Second {
		return nil, errors.New("Get Prometheus block number timeout")
	}

	if model.ValVector != value.Type() {
		return nil, errors.New("Get Prometheus block number data type error")
	}
	nodeBlockNumber := int64(0)
	blockNumber := int64(0)
	for _, data := range value.(model.Vector) {
		number := int64(data.Value)
		if blockNumber <= number {
			blockNumber = number
		}

		instance := string(data.Metric["instance"])
		if !strings.Contains(instance, ":") {
			continue
		}

		ip := strings.Split(instance, ":")[0]
		if ip == nodeIP {
			nodeBlockNumber = number
		}
	}
	blockDiff := blockNumber - nodeBlockNumber

	os, err := getOSVersion(nodeIP)
	if err != nil {
		log.Error("get os version failed", "err", err)
		return nil, err
	}

	gethVersion, err := getGethVersion(nodeIP)
	if err != nil {
		log.Error("get auth node address failed", "err", err)
		return nil, err
	}

	sysTime, err := getNodeSysTime(client, nodeIP)
	if err != nil {
		log.Error("get node sysTime failed", "err", err)
		return nil, err
	}

	upTime, err := getNodeUPTime(client, nodeIP)
	if err != nil {
		log.Error("get node upTime failed", "err", err)
		return nil, err
	}

	memUsage, err := getNodeMemUsage(client, nodeIP)
	if err != nil {
		log.Error("get node memUsage failed", "err", err)
		return nil, err
	}

	cpuUsage, err := getNodeCpuUsage(client, nodeIP)
	if err != nil {
		log.Error("get node cpuUsage failed", "err", err)
		return nil, err
	}

	diskUsage, err := getNodeDiskUsage(client, nodeIP)
	if err != nil {
		log.Error("get node diskUsage failed", "err", err)
		return nil, err
	}
	return &authNodeDetail{BaseInfo: authNodeBaseInfo{IP: nodeIP, OS: os, GethVersion: gethVersion, SysTime: sysTime, RunningTime: upTime, MemUsage: memUsage, CPUUsage: cpuUsage, DiskUsage: diskUsage},
		LinkchainInfo: authNodeChainInfo{nodeBlockNumber, blockDiff}}, nil
}

func getPrometheusAuthNodeCount(client prometheusClient.API) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Cfg.timeout*time.Second)
	defer cancel()
	authNodeCount := 0
	start := time.Now()
	value, err := client.Query(ctx, "latest_block_height", start)
	if err != nil {
		return authNodeCount, err
	}
	end := time.Now()
	if end.Sub(start) > Cfg.timeout*time.Second {
		return authNodeCount, errors.New("Get Prometheus auth node number timeout")
	}

	if model.ValVector != value.Type() {
		return authNodeCount, errors.New("Get Prometheus auth node data type error")
	}

	for _, data := range value.(model.Vector) {
		time := int64(data.Value)
		if time > 0 {
			authNodeCount += 1
		}
	}
	return authNodeCount, nil
}
