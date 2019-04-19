package loader

import (
	"github.com/luyaops/fw/common/log"
	"github.com/luyaops/fw/core"
	"net/url"
	"strings"
)

type RuleStore map[string]core.MethodWrapper

type PrecisionSet []*core.MatchedMethod

func (pq *PrecisionSet) Max() *core.MatchedMethod {
	if pq == nil || *pq == nil {
		return nil
	}
	max := new(core.MatchedMethod)
	for i := 0; i < len(*pq); i++ {
		current := (*pq)[i]
		if max.Precision < current.Precision {
			max = current
		}
	}
	return max
}

func (rs RuleStore) Match(key string) *core.MatchedMethod {
	if v, ok := rs[key]; ok {
		return &core.MatchedMethod{MethodWrapper: v}
	}
	ps := new(PrecisionSet)
	paths := strings.Split(key, "/")

	for keyInDef, methodWrapper := range rs {
		partsInDef := strings.Split(keyInDef, "/")
		if len(paths) == len(partsInDef) {
			values := url.Values{}
			precision := 0
			for i := 0; i < len(paths); i++ {
				if strings.HasPrefix(partsInDef[i], "{") {
					key := strings.TrimSuffix(strings.TrimPrefix(partsInDef[i], "{"), "}")
					values[key] = []string{paths[i]}
					precision = precision + 1
				} else if partsInDef[i] == paths[i] {
					precision = precision + 2
				} else {
					goto NextLoop
				}
			}
			method := core.MatchedMethod{Precision: precision, MergeValues: values, MethodWrapper: methodWrapper}
			log.Debug(method)
			*ps = append(*ps, &method)
		}
	NextLoop:
	}
	return ps.Max()
}
