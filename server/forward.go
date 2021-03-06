package server

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/luyaops/fw/common/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func Run(hostBind string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", forward)

	log.Infof("Listening on %v", hostBind)
	log.Fatal(http.ListenAndServe(hostBind, mux))
}

func forward(w http.ResponseWriter, r *http.Request) {
	// 处理夸域请求
	corsHandle(w, r)
	// 假如请求方法为OPTIONS直接返回204状态码
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	msg, err := handleForward(context.Background(), r)
	if err != nil {
		status, _ := status.FromError(err)
		switch status.Code() {
		case codes.Code(10031):
			http.Redirect(w, r, status.Message(), http.StatusFound)
		default:
			DefaultErrorHandler(w, status.Message(), status.Code())
		}
	} else {
		marshaler := jsonpb.Marshaler{EmitDefaults: true}
		if err := marshaler.Marshal(w, msg); err != nil {
			log.Error(err)
		}
		//w.Write(msg)
	}
}

func corsHandle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Limit,Offset,Origin,Accept,X-Signature,Sec-WebSocket-Protocol")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Access-Control-Max-Age", fmt.Sprint(24*time.Hour/time.Second))

	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
}
