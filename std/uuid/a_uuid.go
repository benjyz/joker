// This file is generated by generate-std.joke script. Do not edit manually!

package uuid

import (
	. "github.com/candid82/joker/core"
)

var __new__P ProcFn = __new_
var new_ Proc = Proc{Fn: __new__P, Name: "new_", Package: "std/joker.uuid"}

func __new_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 0:
		_res := new()
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

func Init() {

	initNative()

	InternsOrThunks()
}

var uuidNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.uuid"))

func init() {
	uuidNamespace.Lazy = Init
}
