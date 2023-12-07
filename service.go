package gnetrpc

import (
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/util"
	"reflect"
	"sync"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfCallMode = reflect.TypeOf((*CallMode)(nil))
var typeOfContext = reflect.TypeOf((*protocol.Context)(nil))
var typeOfSturct = reflect.TypeOf(struct{}{})

type methodType struct {
	sync.Mutex // protects counters
	method     reflect.Method
	ArgType    reflect.Type
	ReplyType  reflect.Type
	// numCalls   uint
}

type functionType struct {
	sync.Mutex // protects counters
	fn         reflect.Value
	ArgType    reflect.Type
	ReplyType  reflect.Type
}

type Service struct {
	name        string                 // name of Service
	value       reflect.Value          // receiver of methods for the Service
	typ         reflect.Type           // type of the receiver
	method      map[string]*methodType // registered methods
	asyncMethod map[string]*methodType
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return util.IsExported(t.Name()) || t.PkgPath() == ""
}

func (s *Service) call(ctx *protocol.Context, mtype *methodType, argv, replyv reflect.Value, isAsync bool) (*CallMode, error) {

	function := mtype.method.Func
	// Invoke the method, providing a new value for the reply.
	var returnValues []reflect.Value
	if isAsync {
		returnValues = function.Call([]reflect.Value{s.value, reflect.ValueOf(ctx), argv, replyv, reflect.ValueOf(struct{}{})})
	} else {
		returnValues = function.Call([]reflect.Value{s.value, reflect.ValueOf(ctx), argv, replyv})
	}
	// The return value for the method is an error.
	callModeInter := returnValues[0].Interface()
	errInter := returnValues[1].Interface()
	var (
		err      error
		callMode *CallMode
	)
	if errInter != nil {
		err = errInter.(error)
	}
	if callModeInter != nil {
		callMode = callModeInter.(*CallMode)
	}
	return callMode, err
}
