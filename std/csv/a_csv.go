// This file is generated by generate-std.joke script. Do not edit manually!

package csv

import (
	. "github.com/candid82/joker/core"
)


var __csv_seq__P ProcFn = __csv_seq_
var csv_seq_ Proc = Proc{Fn: __csv_seq__P, Name: "csv_seq_", Package: "std/csv"}

func __csv_seq_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		rdr := ExtractObject(_args, 0)
		_res := csvSeqOpts(rdr, EmptyArrayMap())
		return _res

	case _c == 2:
		rdr := ExtractObject(_args, 0)
		opts := ExtractMap(_args, 1)
		_res := csvSeqOpts(rdr, opts)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __write__P ProcFn = __write_
var write_ Proc = Proc{Fn: __write__P, Name: "write_", Package: "std/csv"}

func __write_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		f := ExtractIOWriter(_args, 0)
		data := ExtractSeqable(_args, 1)
		_res := write(f, data, EmptyArrayMap())
		return _res

	case _c == 3:
		f := ExtractIOWriter(_args, 0)
		data := ExtractSeqable(_args, 1)
		opts := ExtractMap(_args, 2)
		_res := write(f, data, opts)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __write_string__P ProcFn = __write_string_
var write_string_ Proc = Proc{Fn: __write_string__P, Name: "write_string_", Package: "std/csv"}

func __write_string_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		data := ExtractSeqable(_args, 0)
		_res := writeString(data, EmptyArrayMap())
		return MakeString(_res)

	case _c == 2:
		data := ExtractSeqable(_args, 0)
		opts := ExtractMap(_args, 1)
		_res := writeString(data, opts)
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

var csvNamespace = GLOBAL_ENV.EnsureSymbolIsLib(MakeSymbol("joker.csv"))

func init() {
	csvNamespace.Lazy = Init
}
