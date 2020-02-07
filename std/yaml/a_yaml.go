// This file is generated by generate-std.joke script. Do not edit manually!

package yaml

import (
	. "github.com/candid82/joker/core"
)


var __read_string__P ProcFn = __read_string_
var read_string_ Proc = Proc{Fn: __read_string__P, Name: "read_string_", Package: "std/joker.yaml"}

func __read_string_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := readString(s)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __write_string__P ProcFn = __write_string_
var write_string_ Proc = Proc{Fn: __write_string__P, Name: "write_string_", Package: "std/joker.yaml"}

func __write_string_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		v := ExtractObject(_args, 0)
		_res := writeString(v)
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

var yamlNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.yaml"))

func init() {
	yamlNamespace.Lazy = Init
}
