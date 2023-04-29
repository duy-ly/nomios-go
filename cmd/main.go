package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/duy-ly/nomios-go/consumerpool"
	"github.com/duy-ly/nomios-go/hyperloop"
	"github.com/duy-ly/nomios-go/source"
)

var (
	dbHost       string
	dbPort       int
	dbUser       string
	dbPass       string
	cpCount      int
	cpBufferSize int
	cpFlushTick  string
	streamSize   int
)

func init() {
	flag.StringVar(&dbHost, "db-host", "mysql", "")
	flag.IntVar(&dbPort, "db-port", 3306, "")
	flag.StringVar(&dbUser, "db-user", "nomios", "")
	flag.StringVar(&dbPass, "db-pass", "nomiospass", "")
	flag.IntVar(&cpCount, "cp-count", 5, "")
	flag.IntVar(&cpBufferSize, "cp-buffer-size", 100, "")
	flag.StringVar(&cpFlushTick, "cp-flush-tick", "100ms", "")
	flag.IntVar(&streamSize, "stream-size", 20000, "")
	flag.Parse()
}

func main() {
	cfg := hyperloop.HyperloopConfig{
		SourceConfig: source.SourceConfig{
			Host: dbHost,
			Port: uint16(dbPort),
			User: dbUser,
			Pass: dbPass,

			ServerID: 100,
		},
		PoolConfig: consumerpool.PoolConfig{
			Count:      cpCount,
			BufferSize: cpBufferSize,
		},
		EventStreamSize: streamSize,
	}

	if cpFlushTick != "" {
		if tick, err := time.ParseDuration(cpFlushTick); err == nil {
			cfg.PoolConfig.FlushTick = tick
		} else if n, err := strconv.Atoi(cpFlushTick); err == nil {
			cfg.PoolConfig.FlushTick = time.Duration(n) * time.Millisecond
		}
	}

	h := hyperloop.NewHyperloop(cfg)

	fmt.Println("Nomios starting")

	h.Start()

	fmt.Println("Nomios is started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	for {
		select {
		case <-c:
			fmt.Println("Nomios start graceful terminate")
			h.Stop()
			fmt.Println("Nomios is graceful terminated")
			return
		default:
		}
	}
}
