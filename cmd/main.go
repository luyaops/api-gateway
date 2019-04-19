package main

import (
	"github.com/luyaops/api-gateway/loader"
	"github.com/luyaops/api-gateway/server"
	"github.com/luyaops/fw/common/config"
	"github.com/luyaops/fw/common/log"
)

func main() {
	log.Info("API gateway start")
	loader.Services(config.Endpoints)
	server.Run(config.GatewayBind)
}
