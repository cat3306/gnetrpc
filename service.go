package gnetrpc

import (
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/util"
	"reflect"
	"sync"
)

var typeOfCallMode = reflect.TypeOf((*CallMode)(nil))
var typeOfContext = reflect.TypeOf((*protocol.Context)(nil))
var typeOfSturct = reflect.TypeOf(struct{}{})

type methodType struct {
	sync.Mutex
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

//type functionType struct {
//	sync.Mutex
//	fn         reflect.Value
//	ArgType    reflect.Type
//	ReplyType  reflect.Type
//}

type Service struct {
	name        string
	value       reflect.Value
	typ         reflect.Type
	method      map[string]*methodType
	asyncMethod map[string]*methodType
	preHandler  []Handler
	handlerSet  *HandlerSet
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return util.IsExported(t.Name()) || t.PkgPath() == ""
}

func (s *Service) call(ctx *protocol.Context, mtype *methodType, argv, replyv reflect.Value, isAsync bool) *CallMode {

	function := mtype.method.Func

	var returnValues []reflect.Value
	if isAsync {
		returnValues = function.Call([]reflect.Value{s.value, reflect.ValueOf(ctx), argv, replyv, reflect.ValueOf(struct{}{})})
	} else {
		returnValues = function.Call([]reflect.Value{s.value, reflect.ValueOf(ctx), argv, replyv})
	}

	callModeInter := returnValues[0].Interface()
	var (
		callMode *CallMode
	)
	if callModeInter != nil {
		callMode = callModeInter.(*CallMode)
	}
	return callMode
}
