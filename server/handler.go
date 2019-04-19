package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/luyaops/api-gateway/loader"
	"github.com/luyaops/fw/common/constants"
	"github.com/luyaops/fw/common/log"
	"github.com/luyaops/fw/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func handleForward(ctx context.Context, req *http.Request) (proto.Message, error) {
	body, err := ioutil.ReadAll(req.Body)
	// 获取请求体
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		log.Debug("raw body:", string(body))
		return nil, err
	}
	// 从存储的URL中获取对应的信息
	fmt.Println(req.Method)
	fmt.Println(req.URL.Path)
	sm, err := searchMethod(req.Method, req.URL.Path)
	if err != nil {
		return nil, err
	}
	fmt.Println("GetInputType", sm.Method.GetInputType())
	fmt.Println("GetOutputType", sm.Method.GetOutputType())
	in := protoMessage(sm.Method.GetInputType())
	out := protoMessage(sm.Method.GetOutputType())
	//out := json.RawMessage{}

	json := mergeToBody(string(body), sm.MergeValues, req, in)
	log.Debug(json)

	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err = unmarshaler.Unmarshal(strings.NewReader(json), in); err != nil {
		log.Error(err)
	}
	//sm.package represent for service name by default
	endpoint := sm.Package + ":" + constants.RpcServerPort
	fmt.Println("endpoint:", endpoint)
	fullMethods := "/" + sm.Package + "." + sm.Service + "/" + *sm.Method.Name
	fmt.Println("fullMethod:", fullMethods)
	rpcConn, err := grpc.Dial(endpoint, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rpcConn.Close()
	fullMethod := "/" + sm.Package + "." + sm.Service + "/" + *sm.Method.Name
	fmt.Println("fullMethod:", fullMethod)
	//callOptions := []grpc.CallOption{grpc.CallContentSubtype("mux"), grpc.WaitForReady(false), grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)}
	//if err = rpcConn.Invoke(attachMD(ctx, req, sm.Options), fullMethod, in, out); err != nil {
	if err = rpcConn.Invoke(ctx, fullMethod, in, out); err != nil {
		return nil, err
	}
	return out, nil
}

func searchMethod(method, path string) (*core.MatchedMethod, error) {
	key := method + ":" + path
	matchedMethod := loader.Store.Match(key)
	if matchedMethod != nil {
		return matchedMethod, nil
	}
	return nil, status.Error(codes.NotFound, key+" not been found.")
}

func protoMessage(messageTypeName string) proto.Message {
	messageType := proto.MessageType(messageTypeName)
	return reflect.New(messageType.Elem()).Interface().(proto.Message)
}
