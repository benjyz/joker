// This file is generated by generate-std.joke script. Do not edit manually!

package http

import (
	. "github.com/candid82/joker/core"
)

var __send__P ProcFn = __send_
var send_ Proc = Proc{Fn: __send__P, Name: "send_", Package: "std/http"}

func __send_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		request := ExtractMap(_args, 0)
		_res := sendRequest(request)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __start_file_server__P ProcFn = __start_file_server_
var start_file_server_ Proc = Proc{Fn: __start_file_server__P, Name: "start_file_server_", Package: "std/http"}

func __start_file_server_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		addr := ExtractString(_args, 0)
		root := ExtractString(_args, 1)
		_res := startFileServer(addr, root)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __start_server__P ProcFn = __start_server_
var start_server_ Proc = Proc{Fn: __start_server__P, Name: "start_server_", Package: "std/http"}

func __start_server_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		addr := ExtractString(_args, 0)
		handler := ExtractCallable(_args, 1)
		_res := startServer(addr, handler)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

func Init() {

	initNative()

	InternsOrThunks()
}

var httpNamespace = GLOBAL_ENV.EnsureLib(MakeSymbol("joker.http"))

func init() {
	httpNamespace.Lazy = Init
}
