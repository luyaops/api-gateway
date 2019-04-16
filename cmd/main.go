package main

import (
	"flag"
	"luyaops/api-gateway/loader"
	"luyaops/api-gateway/server"
	"luyaops/fw/common/log"
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
	flag.StringVar(&logLevel, "logLevel", "debug", "Log Level")
	flag.BoolVar(&isHelp, "help", false, "Print this help")
	flag.Parse()

	log.SetLevel(logLevel)
	if isHelp {
		showHelp()
	}
}

func showHelp() {
	flag.PrintDefaults()
	os.Exit(1)
}
