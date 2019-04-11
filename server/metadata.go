package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"luyaops/fw/common/constants"
	"luyaops/fw/common/log"
	"luyaops/fw/core/types"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func attachMD(ctx context.Context, req *http.Request, options map[string]interface{}) context.Context {
	m := make(map[string]string)
	//attach options
	if v, err := json.Marshal(options); err == nil {
		m[constants.RpcMethodOptions] = string(v)
	} else {
		log.Error(err)
	}

	//attach http request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err)
	}

	req.ParseForm()
	newRequest := types.HttpRequest{
		Header: req.Header,
		Method: req.Method,
		Host:   req.Host,
		URL:    req.URL,
		Form:   req.Form,
		Body:   string(body),
	}

	if v, err := json.Marshal(newRequest); err == nil {
		m[constants.HttpRequest] = string(v)
	} else {
		log.Error(err)
	}

	return metadata.NewOutgoingContext(ctx, metadata.New(m))
}
