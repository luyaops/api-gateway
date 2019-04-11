package main

import (
	"flag"
	"luyaops/api-gateway/loader"
	"luyaops/api-gateway/server"
	"luyaops/fw/common/log"
	"os"
)

var hostBind, logLevel string

func main() {
	log.Info("API-Gateway start...")
	loader.Services()
	server.Run(hostBind)
}

func init() {
	var isHelp bool
	flag.StringVar(&hostBind, "bind", ":8080", "Bind address")
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
