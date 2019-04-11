package types

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"net/url"
)

type MethodWrapper struct {
	Package string
	Service string
	Method  *descriptor.MethodDescriptorProto
	Pattern Pattern
	Options map[string]interface{}
}

type Pattern struct {
	Verb string
	Path string
	Body string
}

type MatchedMethod struct {
	MethodWrapper
	Precision   int
	MergeValues url.Values
}
