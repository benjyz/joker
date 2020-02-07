// This file is generated by generate-std.joke script. Do not edit manually!

package hex

import (
	. "github.com/candid82/joker/core"
	"encoding/hex"
)


var __decode_string__P ProcFn = __decode_string_
var decode_string_ Proc = Proc{Fn: __decode_string__P, Name: "decode_string_", Package: "std/joker.hex"}

func __decode_string_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		 t, err := hex.DecodeString(s)
		PanicOnErr(err)
		_res := string(t)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __encode_string__P ProcFn = __encode_string_
var encode_string_ Proc = Proc{Fn: __encode_string__P, Name: "encode_string_", Package: "std/joker.hex"}

func __encode_string_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := hex.EncodeToString([]byte(s))
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

var hexNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.hex"))

func init() {
	hexNamespace.Lazy = Init
}
