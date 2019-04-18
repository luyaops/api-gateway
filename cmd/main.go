package main

import (
	"flag"
	"github.com/luyaops/api-gateway/loader"
	"github.com/luyaops/api-gateway/server"
	"github.com/luyaops/fw/common/log"
	"os"
)

var hostBind, logLevel, etcdAddr string

func main() {
	log.Info("API-Gateway Start...")
	loader.Services(etcdAddr)
	server.Run(hostBind)
}

func init() {
	var isHelp bool
	flag.StringVar(&hostBind, "bind", ":8080", "Bind Address")
	flag.StringVar(&etcdAddr, "etcdAddr", "localhost:2379", "Etcd Address")
	flag.StringVar(&logLevel, "logLevel", "debug", "Log Level: debug info warn error fatal")
	flag.BoolVar(&isHelp, "help", false, "Print this help")
	flag.Parse()

	log.SetLoggerLevel(logLevel)
	if isHelp {
		showHelp()
	}
}

func showHelp() {
	flag.PrintDefaults()
	os.Exit(1)
}
