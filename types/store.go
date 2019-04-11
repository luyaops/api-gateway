package types

import (
	"luyaops/fw/common/log"
	"net/url"
	"strings"
)

type RuleStore map[string]MethodWrapper

func (rs RuleStore) Match(key string) *MatchedMethod {
	if v, ok := rs[key]; ok {
		return &MatchedMethod{MethodWrapper: v}
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
			method := MatchedMethod{Precision: precision, MergeValues: values, MethodWrapper: methodWrapper}
			log.Debug(method)
			*ps = append(*ps, &method)
		}
	NextLoop:
	}
	return ps.Max()
}
