// This file is generated by generate-std.joke script. Do not edit manually!

package math

import (
	. "github.com/candid82/joker/core"
	"math"
)

var e_ Double
var ln_of_10_ Double
var ln_of_2_ Double
var log_10_of_e_ Double
var log_2_of_e_ Double
var max_double_ Double
var phi_ Double
var pi_ Double
var smallest_nonzero_double_ Double
var sqrt_of_2_ Double
var sqrt_of_e_ Double
var sqrt_of_phi_ Double
var sqrt_of_pi_ Double
var __abs__P ProcFn = __abs_
var abs_ Proc = Proc{Fn: __abs__P, Name: "abs_", Package: "std/math"}

func __abs_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Abs(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __ceil__P ProcFn = __ceil_
var ceil_ Proc = Proc{Fn: __ceil__P, Name: "ceil_", Package: "std/math"}

func __ceil_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Ceil(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __copy_sign__P ProcFn = __copy_sign_
var copy_sign_ Proc = Proc{Fn: __copy_sign__P, Name: "copy_sign_", Package: "std/math"}

func __copy_sign_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		x := ExtractNumber(_args, 0)
		y := ExtractNumber(_args, 1)
		_res := math.Copysign(x.Double().D, y.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __cos__P ProcFn = __cos_
var cos_ Proc = Proc{Fn: __cos__P, Name: "cos_", Package: "std/math"}

func __cos_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Cos(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __cube_root__P ProcFn = __cube_root_
var cube_root_ Proc = Proc{Fn: __cube_root__P, Name: "cube_root_", Package: "std/math"}

func __cube_root_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Cbrt(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __dim__P ProcFn = __dim_
var dim_ Proc = Proc{Fn: __dim__P, Name: "dim_", Package: "std/math"}

func __dim_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		x := ExtractNumber(_args, 0)
		y := ExtractNumber(_args, 1)
		_res := math.Dim(x.Double().D, y.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __exp__P ProcFn = __exp_
var exp_ Proc = Proc{Fn: __exp__P, Name: "exp_", Package: "std/math"}

func __exp_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Exp(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __exp_2__P ProcFn = __exp_2_
var exp_2_ Proc = Proc{Fn: __exp_2__P, Name: "exp_2_", Package: "std/math"}

func __exp_2_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Exp2(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __exp_minus_1__P ProcFn = __exp_minus_1_
var exp_minus_1_ Proc = Proc{Fn: __exp_minus_1__P, Name: "exp_minus_1_", Package: "std/math"}

func __exp_minus_1_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Expm1(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __floor__P ProcFn = __floor_
var floor_ Proc = Proc{Fn: __floor__P, Name: "floor_", Package: "std/math"}

func __floor_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Floor(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __hypot__P ProcFn = __hypot_
var hypot_ Proc = Proc{Fn: __hypot__P, Name: "hypot_", Package: "std/math"}

func __hypot_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		p := ExtractNumber(_args, 0)
		q := ExtractNumber(_args, 1)
		_res := math.Hypot(p.Double().D, q.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __inf__P ProcFn = __inf_
var inf_ Proc = Proc{Fn: __inf__P, Name: "inf_", Package: "std/math"}

func __inf_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		sign := ExtractInt(_args, 0)
		_res := math.Inf(sign)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __isinf__P ProcFn = __isinf_
var isinf_ Proc = Proc{Fn: __isinf__P, Name: "isinf_", Package: "std/math"}

func __isinf_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		x := ExtractNumber(_args, 0)
		sign := ExtractInt(_args, 1)
		_res := math.IsInf(x.Double().D, sign)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __log__P ProcFn = __log_
var log_ Proc = Proc{Fn: __log__P, Name: "log_", Package: "std/math"}

func __log_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Log(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __log_10__P ProcFn = __log_10_
var log_10_ Proc = Proc{Fn: __log_10__P, Name: "log_10_", Package: "std/math"}

func __log_10_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Log10(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __log_2__P ProcFn = __log_2_
var log_2_ Proc = Proc{Fn: __log_2__P, Name: "log_2_", Package: "std/math"}

func __log_2_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Log2(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __log_binary__P ProcFn = __log_binary_
var log_binary_ Proc = Proc{Fn: __log_binary__P, Name: "log_binary_", Package: "std/math"}

func __log_binary_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Logb(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __log_plus_1__P ProcFn = __log_plus_1_
var log_plus_1_ Proc = Proc{Fn: __log_plus_1__P, Name: "log_plus_1_", Package: "std/math"}

func __log_plus_1_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Log1p(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __modf__P ProcFn = __modf_
var modf_ Proc = Proc{Fn: __modf__P, Name: "modf_", Package: "std/math"}

func __modf_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := modf(x.Double().D)
		return _res

	default:
		PanicArity(_c)
	}
	return NIL
}

var __nan__P ProcFn = __nan_
var nan_ Proc = Proc{Fn: __nan__P, Name: "nan_", Package: "std/math"}

func __nan_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 0:
		_res := math.NaN()
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __isnan__P ProcFn = __isnan_
var isnan_ Proc = Proc{Fn: __isnan__P, Name: "isnan_", Package: "std/math"}

func __isnan_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.IsNaN(x.Double().D)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __next_after__P ProcFn = __next_after_
var next_after_ Proc = Proc{Fn: __next_after__P, Name: "next_after_", Package: "std/math"}

func __next_after_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		x := ExtractNumber(_args, 0)
		y := ExtractNumber(_args, 1)
		_res := math.Nextafter(x.Double().D, y.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __pow__P ProcFn = __pow_
var pow_ Proc = Proc{Fn: __pow__P, Name: "pow_", Package: "std/math"}

func __pow_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 2:
		x := ExtractNumber(_args, 0)
		y := ExtractNumber(_args, 1)
		_res := math.Pow(x.Double().D, y.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __pow_10__P ProcFn = __pow_10_
var pow_10_ Proc = Proc{Fn: __pow_10__P, Name: "pow_10_", Package: "std/math"}

func __pow_10_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractInt(_args, 0)
		_res := math.Pow10(x)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __round__P ProcFn = __round_
var round_ Proc = Proc{Fn: __round__P, Name: "round_", Package: "std/math"}

func __round_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Round(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __round_to_even__P ProcFn = __round_to_even_
var round_to_even_ Proc = Proc{Fn: __round_to_even__P, Name: "round_to_even_", Package: "std/math"}

func __round_to_even_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.RoundToEven(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __sign_bit__P ProcFn = __sign_bit_
var sign_bit_ Proc = Proc{Fn: __sign_bit__P, Name: "sign_bit_", Package: "std/math"}

func __sign_bit_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Signbit(x.Double().D)
		return MakeBoolean(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __sin__P ProcFn = __sin_
var sin_ Proc = Proc{Fn: __sin__P, Name: "sin_", Package: "std/math"}

func __sin_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Sin(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __sqrt__P ProcFn = __sqrt_
var sqrt_ Proc = Proc{Fn: __sqrt__P, Name: "sqrt_", Package: "std/math"}

func __sqrt_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Sqrt(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

var __trunc__P ProcFn = __trunc_
var trunc_ Proc = Proc{Fn: __trunc__P, Name: "trunc_", Package: "std/math"}

func __trunc_(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 1:
		x := ExtractNumber(_args, 0)
		_res := math.Trunc(x.Double().D)
		return MakeDouble(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

func Init() {
	e_ = MakeDouble(math.E)
	ln_of_10_ = MakeDouble(math.Ln10)
	ln_of_2_ = MakeDouble(math.Ln2)
	log_10_of_e_ = MakeDouble(math.Log10E)
	log_2_of_e_ = MakeDouble(math.Log2E)
	max_double_ = MakeDouble(math.MaxFloat64)
	phi_ = MakeDouble(math.Phi)
	pi_ = MakeDouble(math.Pi)
	smallest_nonzero_double_ = MakeDouble(math.SmallestNonzeroFloat64)
	sqrt_of_2_ = MakeDouble(math.Sqrt2)
	sqrt_of_e_ = MakeDouble(math.SqrtE)
	sqrt_of_phi_ = MakeDouble(math.SqrtPhi)
	sqrt_of_pi_ = MakeDouble(math.SqrtPi)
	initNative()

	InternsOrThunks()
}

var mathNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.math"))

func init() {
	mathNamespace.Lazy = Init
}
