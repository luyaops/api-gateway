package server

import (
	"bytes"
	"context"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"luyaops/api-gateway/loader"
	"luyaops/api-gateway/types"
	"luyaops/fw/common/constants"
	"luyaops/fw/common/log"
	"net/http"
	"reflect"
	"strings"
)

func handleForward(ctx context.Context, req *http.Request) (proto.Message, error) {
	body, err := ioutil.ReadAll(req.Body)
	//body will be consumed again
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		log.Debug("raw body:", string(body))
		return nil, err
	}

	sm, err := searchMethod(req.Method, req.URL.Path)
	if err != nil {
		return nil, err
	}

	in := protoMessage(sm.Method.GetInputType())
	out := protoMessage(sm.Method.GetOutputType())

	json := mergeToBody(string(body), sm.MergeValues, req, in)
	log.Debug(json)

	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err = unmarshaler.Unmarshal(strings.NewReader(json), in); err != nil {
		log.Error(err)
	}
	//sm.package represent for service name by default
	endpoint := sm.Package + ":" + constants.RpcServerPort
	rpcConn, err := grpc.Dial(endpoint, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rpcConn.Close()
	fullMethod := "/" + sm.Package + "." + sm.Service + "/" + *sm.Method.Name
	//if err = rpcConn.Invoke(attachMD(ctx, req, sm.Options), fullMethod, in, out); err != nil {
	if err = rpcConn.Invoke(ctx, fullMethod, in, out); err != nil {
		return nil, err
	}
	return out, nil
}

func searchMethod(method, path string) (*types.MatchedMethod, error) {
	key := method + ":" + path
	matchedMethod := loader.RuleStore.Match(key)
	if matchedMethod != nil {
		return matchedMethod, nil
	}
	return nil, status.Error(codes.NotFound, key+" not been found.")
}

func protoMessage(messageTypeName string) proto.Message {
	typeName := strings.TrimLeft(messageTypeName, ".")
	messageType := proto.MessageType(typeName)
	return reflect.New(messageType.Elem()).Interface().(proto.Message)
}
