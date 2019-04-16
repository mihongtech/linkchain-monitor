package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mihongtech/linkchain/common/util/log"
	"gopkg.in/urfave/cli.v1"
)

// main is just a boring entry point to set up the CLI app.
func main() {
	app := cli.NewApp()
	app.Name = "explorer"
	app.Usage = "linkchain explorer server"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 8090,
			Usage: "http service linsten port",
		},

		cli.IntFlag{
			Name:  "loglevel",
			Value: 3,
			Usage: "log level to emit to the screen",
		},

		cli.IntFlag{
			Name:  "prometheusport",
			Value: 9090,
			Usage: "prometheus server port",
		},

		cli.IntFlag{
			Name:  "rpcport",
			Value: 8888,
			Usage: "geth rpc port",
		},

		cli.IntFlag{
			Name:  "timeout",
			Value: 10,
			Usage: "http timeout",
		},
	}

	app.Action = func(c *cli.Context) error {
		// Set up the logger to print everything and the random generator
		log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(c.Int("loglevel")), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))
		rand.Seed(time.Now().UnixNano())
		Cfg = &config{
			ip:      "127.0.0.1",
			port:    c.Int("prometheusport"),
			rpcPort: c.Int("rpcport"),
			timeout: time.Duration(c.Int("timeout")),
		}

		s := &Statistics{
			RpcPort:  c.Int("rpcport"),
			Internal: c.Int("interval"),
		}
		go s.Serve()

		router := NewRouter()
		log.Debug("start node ctl web service", "log", http.ListenAndServe(":"+strconv.Itoa(c.Int("port")), router))

		return nil
	}
	app.Run(os.Args)
}
