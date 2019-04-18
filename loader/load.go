package loader

import (
	"encoding/json"
	"github.com/luyaops/api-gateway/types"
	_ "github.com/luyaops/example/proto"
	"github.com/luyaops/fw/common/constants"
	"github.com/luyaops/fw/common/etcd"
	"github.com/luyaops/fw/common/log"
	"go.etcd.io/etcd/clientv3"
)

var RuleStore = make(types.RuleStore)

func Services(endpoints []string) {
	load(endpoints)
}

const (
	API = "api"
)

func load(endpoints []string) {
	var methods []types.MethodWrapper
	client := etcd.NewStore(endpoints)
	if registration, err := client.Get(constants.RegistryPrefix, clientv3.WithPrefix()); err != nil {
		log.Fatalf("Failed to load registration information:%v", err)
	} else {
		for _, v := range registration.Kvs {
			err := json.Unmarshal(v.Value, &methods)
			if err != nil {
				log.Error(err)
			}
		}
	}
	for _, md := range methods {
		//key := md.Pattern.Verb + ":" + md.Pattern.Path
		key := md.Pattern.Verb + ":/" + API + md.Pattern.Path
		log.Debug(key, " --> ", md.Service, ".", *md.Method.Name)
		RuleStore[key] = md
	}
	log.Infof("Registration information loaded successfully")
}
