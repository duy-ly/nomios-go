package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/duy-ly/nomios-go/hyperloop"
	"github.com/duy-ly/nomios-go/logger"
	"github.com/spf13/viper"
)

var (
	confFile     string
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
	flag.StringVar(&confFile, "conf-file", "", "")

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
	err := loadConfig()
	if err != nil {
		logger.NomiosLog.Panic("Nomios cannot start due to unable to load config ", err)
	}

	h := hyperloop.NewHyperloop()

	logger.NomiosLog.Debug("Nomios starting.")

	h.Start()

	logger.NomiosLog.Debug("Nomios is started.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	for {
		select {
		case <-c:
			logger.NomiosLog.Debug("Nomios start graceful terminate.")
			h.Stop()
			logger.NomiosLog.Debug("Nomios is graceful terminated.")
			return
		default:
		}
	}
}

func loadConfig() error {
	viper.SetConfigName("nomios")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}
