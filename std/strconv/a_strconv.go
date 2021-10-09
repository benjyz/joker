// This file is generated by generate-std.joke script. Do not edit manually!

package strconv

import (
	. "github.com/candid82/joker/core"
	"strconv"
)


var __atoi__P ProcFn = __atoi_
var atoi_ Proc = Proc{Fn: __atoi__P, Name: "atoi_", Package: "std/strconv"}

func __atoi_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res, err := strconv.Atoi(s)
		PanicOnErr(err)
		return MakeInt(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __iscan_backquote__P ProcFn = __iscan_backquote_
var iscan_backquote_ Proc = Proc{Fn: __iscan_backquote__P, Name: "iscan_backquote_", Package: "std/strconv"}

func __iscan_backquote_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := strconv.CanBackquote(s)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __format_bool__P ProcFn = __format_bool_
var format_bool_ Proc = Proc{Fn: __format_bool__P, Name: "format_bool_", Package: "std/strconv"}

func __format_bool_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		b := ExtractBoolean(_args, 0)
		_res := strconv.FormatBool(b)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __format_double__P ProcFn = __format_double_
var format_double_ Proc = Proc{Fn: __format_double__P, Name: "format_double_", Package: "std/strconv"}

func __format_double_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 4:
		f := ExtractDouble(_args, 0)
		fmt := ExtractChar(_args, 1)
		prec := ExtractInt(_args, 2)
		bitSize := ExtractInt(_args, 3)
		_res := strconv.FormatFloat(f, byte(fmt), prec, bitSize)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __format_int__P ProcFn = __format_int_
var format_int_ Proc = Proc{Fn: __format_int__P, Name: "format_int_", Package: "std/strconv"}

func __format_int_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		i := ExtractInt(_args, 0)
		base := ExtractInt(_args, 1)
		_res := strconv.FormatInt(int64(i), base)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __isgraphic__P ProcFn = __isgraphic_
var isgraphic_ Proc = Proc{Fn: __isgraphic__P, Name: "isgraphic_", Package: "std/strconv"}

func __isgraphic_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		c := ExtractChar(_args, 0)
		_res := strconv.IsGraphic(c)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __itoa__P ProcFn = __itoa_
var itoa_ Proc = Proc{Fn: __itoa__P, Name: "itoa_", Package: "std/strconv"}

func __itoa_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		i := ExtractInt(_args, 0)
		_res := strconv.Itoa(i)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __parse_bool__P ProcFn = __parse_bool_
var parse_bool_ Proc = Proc{Fn: __parse_bool__P, Name: "parse_bool_", Package: "std/strconv"}

func __parse_bool_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res, err := strconv.ParseBool(s)
		PanicOnErr(err)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __parse_double__P ProcFn = __parse_double_
var parse_double_ Proc = Proc{Fn: __parse_double__P, Name: "parse_double_", Package: "std/strconv"}

func __parse_double_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res, err := strconv.ParseFloat(s, 64)
		PanicOnErr(err)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __parse_int__P ProcFn = __parse_int_
var parse_int_ Proc = Proc{Fn: __parse_int__P, Name: "parse_int_", Package: "std/strconv"}

func __parse_int_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 3:
		s := ExtractString(_args, 0)
		base := ExtractInt(_args, 1)
		bitSize := ExtractInt(_args, 2)
		t, err := strconv.ParseInt(s, base, bitSize)
		PanicOnErr(err)
		_res := int(t)
		return MakeInt(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __isprintable__P ProcFn = __isprintable_
var isprintable_ Proc = Proc{Fn: __isprintable__P, Name: "isprintable_", Package: "std/strconv"}

func __isprintable_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		c := ExtractChar(_args, 0)
		_res := strconv.IsPrint(c)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __quote__P ProcFn = __quote_
var quote_ Proc = Proc{Fn: __quote__P, Name: "quote_", Package: "std/strconv"}

func __quote_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := strconv.Quote(s)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __quote_char__P ProcFn = __quote_char_
var quote_char_ Proc = Proc{Fn: __quote_char__P, Name: "quote_char_", Package: "std/strconv"}

func __quote_char_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		c := ExtractChar(_args, 0)
		_res := strconv.QuoteRune(c)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __quote_char_to_ascii__P ProcFn = __quote_char_to_ascii_
var quote_char_to_ascii_ Proc = Proc{Fn: __quote_char_to_ascii__P, Name: "quote_char_to_ascii_", Package: "std/strconv"}

func __quote_char_to_ascii_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		c := ExtractChar(_args, 0)
		_res := strconv.QuoteRuneToASCII(c)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __quote_char_to_graphic__P ProcFn = __quote_char_to_graphic_
var quote_char_to_graphic_ Proc = Proc{Fn: __quote_char_to_graphic__P, Name: "quote_char_to_graphic_", Package: "std/strconv"}

func __quote_char_to_graphic_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		c := ExtractChar(_args, 0)
		_res := strconv.QuoteRuneToGraphic(c)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __quote_to_ascii__P ProcFn = __quote_to_ascii_
var quote_to_ascii_ Proc = Proc{Fn: __quote_to_ascii__P, Name: "quote_to_ascii_", Package: "std/strconv"}

func __quote_to_ascii_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := strconv.QuoteToASCII(s)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __quote_to_graphic__P ProcFn = __quote_to_graphic_
var quote_to_graphic_ Proc = Proc{Fn: __quote_to_graphic__P, Name: "quote_to_graphic_", Package: "std/strconv"}

func __quote_to_graphic_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := strconv.QuoteToGraphic(s)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __unquote__P ProcFn = __unquote_
var unquote_ Proc = Proc{Fn: __unquote__P, Name: "unquote_", Package: "std/strconv"}

func __unquote_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res, err := strconv.Unquote(s)
		PanicOnErr(err)
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

var strconvNamespace = GLOBAL_ENV.EnsureSymbolIsLib(MakeSymbol("joker.strconv"))

func init() {
	strconvNamespace.Lazy = Init
}
