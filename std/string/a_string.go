// This file is generated by generate-std.joke script. Do not edit manually!

package string

import (
	. "github.com/candid82/joker/core"
	"regexp"
	"strings"
	"unicode"
)

var __isblank__P ProcFn = __isblank_
var isblank_ Proc = Proc{Fn: __isblank__P, Name: "isblank_", Package: "std/joker.string"}

func __isblank_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractObject(_args, 0)
		_res := isBlank(s)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __capitalize__P ProcFn = __capitalize_
var capitalize_ Proc = Proc{Fn: __capitalize__P, Name: "capitalize_", Package: "std/joker.string"}

func __capitalize_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractStringable(_args, 0)
		_res := capitalize(s)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __isends_with__P ProcFn = __isends_with_
var isends_with_ Proc = Proc{Fn: __isends_with__P, Name: "isends_with_", Package: "std/joker.string"}

func __isends_with_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		s := ExtractString(_args, 0)
		substr := ExtractStringable(_args, 1)
		_res := strings.HasSuffix(s, substr)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __escape__P ProcFn = __escape_
var escape_ Proc = Proc{Fn: __escape__P, Name: "escape_", Package: "std/joker.string"}

func __escape_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		s := ExtractString(_args, 0)
		cmap := ExtractCallable(_args, 1)
		_res := escape(s, cmap)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __isincludes__P ProcFn = __isincludes_
var isincludes_ Proc = Proc{Fn: __isincludes__P, Name: "isincludes_", Package: "std/joker.string"}

func __isincludes_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		s := ExtractString(_args, 0)
		substr := ExtractStringable(_args, 1)
		_res := strings.Contains(s, substr)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __index_of__P ProcFn = __index_of_
var index_of_ Proc = Proc{Fn: __index_of__P, Name: "index_of_", Package: "std/joker.string"}

func __index_of_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		s := ExtractString(_args, 0)
		value := ExtractObject(_args, 1)
		_res := indexOf(s, value, 0)
		return _res

	case _c == 3:
		s := ExtractString(_args, 0)
		value := ExtractObject(_args, 1)
		from := ExtractInt(_args, 2)
		_res := indexOf(s, value, from)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __join__P ProcFn = __join_
var join_ Proc = Proc{Fn: __join__P, Name: "join_", Package: "std/joker.string"}

func __join_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		coll := ExtractSeqable(_args, 0)
		_res := join("", coll)
		return MakeString(_res)

	case _c == 2:
		separator := ExtractStringable(_args, 0)
		coll := ExtractSeqable(_args, 1)
		_res := join(separator, coll)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __last_index_of__P ProcFn = __last_index_of_
var last_index_of_ Proc = Proc{Fn: __last_index_of__P, Name: "last_index_of_", Package: "std/joker.string"}

func __last_index_of_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		s := ExtractString(_args, 0)
		value := ExtractObject(_args, 1)
		_res := lastIndexOf(s, value, 0)
		return _res

	case _c == 3:
		s := ExtractString(_args, 0)
		value := ExtractObject(_args, 1)
		from := ExtractInt(_args, 2)
		_res := lastIndexOf(s, value, from)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __lower_case__P ProcFn = __lower_case_
var lower_case_ Proc = Proc{Fn: __lower_case__P, Name: "lower_case_", Package: "std/joker.string"}

func __lower_case_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractStringable(_args, 0)
		_res := strings.ToLower(s)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __pad_left__P ProcFn = __pad_left_
var pad_left_ Proc = Proc{Fn: __pad_left__P, Name: "pad_left_", Package: "std/joker.string"}

func __pad_left_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 3:
		s := ExtractString(_args, 0)
		pad := ExtractStringable(_args, 1)
		n := ExtractInt(_args, 2)
		_res := padLeft(s, pad, n)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __pad_right__P ProcFn = __pad_right_
var pad_right_ Proc = Proc{Fn: __pad_right__P, Name: "pad_right_", Package: "std/joker.string"}

func __pad_right_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 3:
		s := ExtractString(_args, 0)
		pad := ExtractStringable(_args, 1)
		n := ExtractInt(_args, 2)
		_res := padRight(s, pad, n)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __re_quote__P ProcFn = __re_quote_
var re_quote_ Proc = Proc{Fn: __re_quote__P, Name: "re_quote_", Package: "std/joker.string"}

func __re_quote_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := regexp.MustCompile(regexp.QuoteMeta(s))
		return MakeRegex(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __replace__P ProcFn = __replace_
var replace_ Proc = Proc{Fn: __replace__P, Name: "replace_", Package: "std/joker.string"}

func __replace_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 3:
		s := ExtractString(_args, 0)
		match := ExtractObject(_args, 1)
		repl := ExtractStringable(_args, 2)
		_res := replace(s, match, repl)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __replace_first__P ProcFn = __replace_first_
var replace_first_ Proc = Proc{Fn: __replace_first__P, Name: "replace_first_", Package: "std/joker.string"}

func __replace_first_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 3:
		s := ExtractString(_args, 0)
		match := ExtractObject(_args, 1)
		repl := ExtractStringable(_args, 2)
		_res := replaceFirst(s, match, repl)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __reverse__P ProcFn = __reverse_
var reverse_ Proc = Proc{Fn: __reverse__P, Name: "reverse_", Package: "std/joker.string"}

func __reverse_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := reverse(s)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __split__P ProcFn = __split_
var split_ Proc = Proc{Fn: __split__P, Name: "split_", Package: "std/joker.string"}

func __split_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		s := ExtractString(_args, 0)
		sep := ExtractObject(_args, 1)
		_res := splitOnStringOrRegex(s, sep, 0)
		return _res

	case _c == 3:
		s := ExtractString(_args, 0)
		sep := ExtractObject(_args, 1)
		n := ExtractInt(_args, 2)
		_res := splitOnStringOrRegex(s, sep, n)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __split_lines__P ProcFn = __split_lines_
var split_lines_ Proc = Proc{Fn: __split_lines__P, Name: "split_lines_", Package: "std/joker.string"}

func __split_lines_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := split(s, newLine, 0)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __isstarts_with__P ProcFn = __isstarts_with_
var isstarts_with_ Proc = Proc{Fn: __isstarts_with__P, Name: "isstarts_with_", Package: "std/joker.string"}

func __isstarts_with_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		s := ExtractString(_args, 0)
		substr := ExtractStringable(_args, 1)
		_res := strings.HasPrefix(s, substr)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __trim__P ProcFn = __trim_
var trim_ Proc = Proc{Fn: __trim__P, Name: "trim_", Package: "std/joker.string"}

func __trim_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := strings.TrimSpace(s)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __trim_left__P ProcFn = __trim_left_
var trim_left_ Proc = Proc{Fn: __trim_left__P, Name: "trim_left_", Package: "std/joker.string"}

func __trim_left_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := strings.TrimLeftFunc(s, unicode.IsSpace)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __trim_newline__P ProcFn = __trim_newline_
var trim_newline_ Proc = Proc{Fn: __trim_newline__P, Name: "trim_newline_", Package: "std/joker.string"}

func __trim_newline_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := strings.TrimRight(s, "\n\r")
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __trim_right__P ProcFn = __trim_right_
var trim_right_ Proc = Proc{Fn: __trim_right__P, Name: "trim_right_", Package: "std/joker.string"}

func __trim_right_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := strings.TrimRightFunc(s, unicode.IsSpace)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __triml__P ProcFn = __triml_
var triml_ Proc = Proc{Fn: __triml__P, Name: "triml_", Package: "std/joker.string"}

func __triml_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := strings.TrimLeftFunc(s, unicode.IsSpace)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __trimr__P ProcFn = __trimr_
var trimr_ Proc = Proc{Fn: __trimr__P, Name: "trimr_", Package: "std/joker.string"}

func __trimr_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractString(_args, 0)
		_res := strings.TrimRightFunc(s, unicode.IsSpace)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __upper_case__P ProcFn = __upper_case_
var upper_case_ Proc = Proc{Fn: __upper_case__P, Name: "upper_case_", Package: "std/joker.string"}

func __upper_case_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		s := ExtractStringable(_args, 0)
		_res := strings.ToUpper(s)
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

var stringNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.string"))

func init() {
	stringNamespace.Lazy = Init
}
