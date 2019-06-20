package loader

import (
	"encoding/json"
	_ "github.com/luyaops/cmdb/proto"
	_ "github.com/luyaops/example/proto"
	"github.com/luyaops/fw/common/constants"
	"github.com/luyaops/fw/common/etcd"
	"github.com/luyaops/fw/common/log"
	"github.com/luyaops/fw/core"
	"go.etcd.io/etcd/clientv3"
)

var Store = make(RuleStore)

func Services(endpoints []string) {
	load(endpoints)
}

const (
	API = "api"
)

func load(endpoints []string) {
	var methods []core.MethodWrapper
	client := etcd.NewStore(endpoints)
	if registration, err := client.Get(constants.RegistryPrefix, clientv3.WithPrefix()); err != nil {
		log.Fatalf("Failed to load registration information:%v", err)
	} else {
		for _, v := range registration.Kvs {
			//fmt.Printf("%+v\n", string(v.Value))
			var method []core.MethodWrapper
			err := json.Unmarshal(v.Value, &method)
			methods = append(methods, method...)
			if err != nil {
				log.Error(err)
			}
		}
	}

	for _, md := range methods {
		key := md.Pattern.Verb + ":/" + API + md.Pattern.Path
		log.Debug(key, " --> ", md.Service, ".", *md.Method.Name)
		Store[key] = md
	}
	log.Infof("Registration information loaded successfully")
}
