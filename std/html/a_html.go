// This file is generated by generate-std.joke script. Do not edit manually!

package html

import (
	. "github.com/candid82/joker/core"
	"html"
)

var __escape__P ProcFn = __escape_
var escape_ Proc = Proc{Fn: __escape__P, Name: "escape_", Package: "std/html"}

func __escape_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := html.EscapeString(s)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __unescape__P ProcFn = __unescape_
var unescape_ Proc = Proc{Fn: __unescape__P, Name: "unescape_", Package: "std/html"}

func __unescape_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := html.UnescapeString(s)
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

var htmlNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.html"))

func init() {
	htmlNamespace.Lazy = Init
}
