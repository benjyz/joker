// This file is generated by generate-std.joke script. Do not edit manually!

package httptest

import (
	. "github.com/candid82/joker/core"
	"net/http/httptest"
)

var httptestNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("go.std.net.http.httptest"))

var DefaultRemoteAddr_ String

var NewRecorder_ Proc

func __NewRecorder_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 0:
		_res := __newRecorder()
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

func Init() {
	DefaultRemoteAddr_ = MakeString(httptest.DefaultRemoteAddr)
	NewRecorder_ = __NewRecorder_

	initNative()

	httptestNamespace.ResetMeta(MakeMeta(nil, `Provides a low-level interface to the net/http/httptest package.

Package httptest provides utilities for HTTP testing.
`, "1.0"))

	httptestNamespace.InternVar("DefaultRemoteAddr", DefaultRemoteAddr_,
		MakeMeta(
			nil,
			`DefaultRemoteAddr is the default remote address to return in RemoteAddr if
an explicit DefaultRemoteAddr isn't set on ResponseRecorder.
`, "1.0"))

	httptestNamespace.InternVar("NewRecorder", NewRecorder_,
		MakeMeta(
			NewListFrom(NewVectorFrom()),
			`NewRecorder returns an initialized ResponseRecorder.

Go return type: *ResponseRecorder

Joker input arguments: []

Joker return type: (atom-of go.std.net.http.httptest/ResponseRecorder)`, "1.0"))

}

func init() {
	httptestNamespace.Lazy = Init
}
