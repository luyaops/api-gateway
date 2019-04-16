package loader

import (
	"encoding/json"
	"luyaops/api-gateway/types"
	_ "luyaops/example/proto"
	"luyaops/fw/common/etcd"
	"luyaops/fw/common/log"
)

var RuleStore = make(types.RuleStore)

func Services(etcdAddr string) {
	load(etcdAddr)
}

const (
	API = "api"
)

func load(etcdAddr string) {
	var methods []types.MethodWrapper
	client := etcd.NewStore(etcdAddr)
	say, _ := client.Get("say")
	for _, v := range say.Kvs {
		err := json.Unmarshal(v.Value, &methods)
		if err != nil {
			log.Error(err)
		}
	}
	//err := json.Unmarshal([]byte(string(PROTO_JSON)), &methods)
	//if err != nil {
	//	log.Error(err)
	//}
	for _, md := range methods {
		//key := md.Pattern.Verb + ":" + md.Pattern.Path
		key := md.Pattern.Verb + ":/" + API + md.Pattern.Path
		log.Debug(key, "-->", md.Service, ".", *md.Method.Name)
		RuleStore[key] = md
	}
}
