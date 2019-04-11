package server

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

func mergeToBody(bodyJSON string, pathValues url.Values, req *http.Request, pbMsg interface{}) string {
	//need firstly call ParseForm before req.Form
	req.ParseForm()
	queryJSON := toJSONStr(pbMsg, req.Form)
	pathJSON := toJSONStr(pbMsg, pathValues)

	if bodyJSON == "" {
		bodyJSON = "{}"
	}

	jsonStr := strings.TrimSuffix(bodyJSON, "}") + pathJSON + queryJSON + "}"
	replacer := strings.NewReplacer("{,", "{")
	return replacer.Replace(jsonStr)
}

func toJSONStr(pbMsg interface{}, values url.Values) (str string) {
	for k, v := range values {
		field := reflect.ValueOf(pbMsg).Elem().FieldByName(strings.Title(k))
		if field.IsValid() {
			switch field.Type().Name() {
			case "int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"float32", "float64", "bool":
				str = str + ",\"" + k + "\":" + v[0] + ""
			default:
				str = str + ",\"" + k + "\":\"" + v[0] + "\""
			}
		}
	}

	return str
}
