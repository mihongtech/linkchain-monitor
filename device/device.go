package main

import (
	_ "bytes"
	"flag"
	"net/http"
	"os"
	_ "os/exec"
	_ "path/filepath"
	"strings"
	"time"

	"github.com/mihongtech/linkchain/common/util/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var (
	sysTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sys_time",
		Help: "sys time",
	})
	upTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "up_time",
		Help: "sys up time",
	})
	memUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mem_usage",
		Help: "mem usage",
	})
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "cpu usage",
	})
	diskUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "disk_usage",
		Help: "disk usage",
	},
		[]string{"device"})
	blockHeight int64 = 0
)

func init() {
	prometheus.MustRegister(sysTime)
	prometheus.MustRegister(upTime)
	prometheus.MustRegister(memUsage)
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(diskUsage)
}

func updateBaseInfo(interval int) {
	sysTime.Set(float64(time.Now().Unix()))
	v, _ := mem.VirtualMemory()
	memUsage.Set(v.UsedPercent / float64(100))
	c, _ := cpu.Percent(time.Duration(interval)*time.Second, false)
	cpuUsage.Set(c[0] / float64(100))
	partitions, _ := disk.Partitions(true)
	for _, p := range partitions {
		// log.Debug("partition is ", "p", p)
		if p.Fstype == "tmpfs" || !strings.Contains(p.Device, "/dev") {
			continue
		}
		d, _ := disk.Usage(p.Mountpoint)
		if d == nil || d.Total == 0 {
			log.Error("error partition is ", "p", p)
			continue
		}
		// log.Info("d vaule is", "d", d)
		diskUsage.With(prometheus.Labels{"device": p.Device}).Set(d.UsedPercent / float64(100))
	}
	u, _ := host.Uptime()
	upTime.Set(float64(u))
}

func main() {
	var (
		listenAddr    = flag.String("addr", ":9200", "listen address")
		fetchInterval = flag.Int("interval", 5, "fetch interval")
		loglevel      = flag.Int("loglevel", 3, "log level")
		metrics       = flag.String("metrics_path", "/device", "the metrics path of promethus")
	)
	flag.Parse()
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(*loglevel), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	go func() {
		for {
			<-time.After(time.Duration(*fetchInterval) * time.Second)
			updateBaseInfo(*fetchInterval)
		}
	}()

	http.Handle(*metrics, prometheus.Handler())
	http.ListenAndServe(*listenAddr, nil)
}
