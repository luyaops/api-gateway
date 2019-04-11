package loader

import (
	"encoding/json"
	"luyaops/api-gateway/types"
	"luyaops/fw/common/log"
)

var RuleStore = make(types.RuleStore)

func Services() {
	load()
}

const (
	API = "api"
)

func load() {
	//log.Debug(PROTO_JSON)
	var methods []types.MethodWrapper
	err := json.Unmarshal([]byte(string(PROTO_JSON)), &methods)
	if err != nil {
		log.Error(err)
	}

	for _, md := range methods {
		//key := md.Pattern.Verb + ":" + md.Pattern.Path
		key := md.Pattern.Verb + ":/" + API + md.Pattern.Path
		log.Debug(key, "-->", md.Service, ".", *md.Method.Name)
		RuleStore[key] = md
	}
}
