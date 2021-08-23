package jtypes

import (
	"fmt"
	"github.com/candid82/joker/tools/gostd/astutils"
	"github.com/candid82/joker/tools/gostd/genutils"
	. "github.com/candid82/joker/tools/gostd/godb"
	. "go/ast"
)

// Info on Clojure types, including map of Clojure type names to said type
// info.  A Clojure type name is either unqualified (built-in, not
// namespace-rooted) or fully qualified by a namespace name
// (e.g. "go.std.example/SomeType").
type Info struct {
	Expr                 Expr   // [key] The canonical referencing expression (if any)
	FullName             string // [key] Full name of type as a Clojure expression
	FullNameDoc          string // Full name of type as a Clojure expression (for documentation); just e.g. "Int" for builtin global types
	Who                  string // who made me
	Pattern              string // E.g. "%s", "refTo%s" (for reference types), "arrayOf%s" (for array types)
	Namespace            string // E.g. "go.std.net.url", in which this type resides ("" denotes global namespace)
	BaseName             string // E.g. "Listener"
	BaseNameDoc          string // Might be e.g. "Int" when BaseName is "Int8"
	ArgClojureType       string // Can convert this type to a Go function arg with my type
	ArgFromClojureObject string // Append this to Clojure object to extract value of my type
	ArgExtractFunc       string
	ArgClojureArgType    string // Clojure argument type for a Go function arg with my type
	ConvertFromClojure   string // Pattern to convert a (scalar) %s to this type
	ConvertFromMap       string // Pattern to convert a map %s key %s to this type
	AsClojureObject      string // Pattern to convert this type to a normal Clojure type; empty string means wrap in a GoObject
	PromoteType          string // Pattern to promote to a canonical type (used by constant evaluation)
	GoApiString          string // E.g. "error", "uint16", "go.std.net/IPv4Addr"
	AliasOf              **Info // What this aliases to (nil if nothing)
	IsUnsupported        bool   // Is this unsupported?
}

// Maps type-defining Expr or Clojure type names (with or without
// "<ns>/" prefixes, depending on globality) to exactly one struct
// describing that type.
var typesByExpr = map[Expr]*Info{}
var typesByFullname = map[string]*Info{}

func patternForExpr(e Expr) (pattern string, ue Expr) {
	switch v := e.(type) {
	case *ArrayType:
		len, _ := astutils.IntExprToString(v.Len)
		pattern, e = patternForExpr(v.Elt)
		return "array" + len + "Of" + pattern, e
	case *StarExpr:
		pattern, e = patternForExpr(v.X)
		return "refTo" + pattern, e
	case *MapType:
		patternKey, _ := patternForExpr(v.Key)
		patternValue, eValue := patternForExpr(v.Value)
		res := "map_" + patternKey + "_Of_" + fmt.Sprintf(patternValue, "<whatever>")
		return fmt.Sprintf("ABEND777(jtypes.go: multiple underlying expressions not supported: %s)", res), eValue
	case *ChanType:
		pattern, e = patternForExpr(v.Value)
		baseName := "chan"
		switch v.Dir {
		case SEND:
			baseName = "chanSend"
		case RECV:
			baseName = "chanRecv"
		case SEND | RECV:
		default:
			baseName = fmt.Sprintf("ABEND737(jtypes.go: %s Dir=0x%x not supported)", astutils.ExprToString(v), v.Dir)
		}
		return baseName + "Of" + pattern, e
	default:
		return "%s", e
	}
}

const resolveAlias = false

func namingForExpr(e Expr) (pattern, ns, baseName, baseNameDoc, name, nameDoc, goApiString string, info *Info) {
	var ue Expr
	pattern, ue = patternForExpr(e)

	switch v := ue.(type) {
	case *Ident:
		if !astutils.IsBuiltin(v.Name) {
			ns = ClojureNamespaceForExpr(ue)
			baseName = v.Name
			baseNameDoc = baseName
			goApiString = baseName
		} else {
			uInfo, found := goTypeMap[v.Name]
			if !found {
				panic(fmt.Sprintf("no type info for builtin `%s'", v.Name))
			}
			if resolveAlias && uInfo.AliasOf != nil {
				uInfo = *uInfo.AliasOf
			}
			baseName = uInfo.FullName
			baseNameDoc = uInfo.FullNameDoc
			goApiString = uInfo.GoApiString
			if e == ue {
				info = uInfo
			}
		}
	case *SelectorExpr:
		pkgName := v.X.(*Ident).Name
		ns = ClojureNamespaceForGoFile(pkgName, GoFileForExpr(v))
		baseName = v.Sel.Name
		baseNameDoc = baseName
		goApiString = baseName
	case *InterfaceType:
		if !v.Incomplete && astutils.IsEmptyFieldList(v.Methods) {
			baseName = "GoObject"
		} else {
			baseName = fmt.Sprintf("ABEND320(jtypes.go: %s not supported)", astutils.ExprToString(v))
		}
		baseNameDoc = baseName
		goApiString = baseName
	case *StructType:
		if astutils.IsEmptyFieldList(v.Fields) {
			baseName = "GoObject"
		} else {
			baseName = fmt.Sprintf("ABEND787(jtypes.go: %s not supported)", astutils.ExprToString(v))
		}
		baseNameDoc = baseName
		goApiString = baseName
	case *FuncType:
		if astutils.IsEmptyFieldList(v.Params) && astutils.IsEmptyFieldList(v.Results) {
			baseName = "func"
		} else {
			baseName = fmt.Sprintf("ABEND727(jtypes.go: %s not supported)", astutils.ExprToString(v))
		}
		baseNameDoc = baseName
		goApiString = baseName
	case *Ellipsis:
		baseName = fmt.Sprintf("ABEND747(jtypes.go: %s not supported)", astutils.ExprToString(v))
		baseNameDoc = baseName
		goApiString = baseName
	default:
		panic(fmt.Sprintf("unrecognized underlying expr %T for %T", ue, e))
	}

	name = genutils.CombineClojureName(ns, fmt.Sprintf(pattern, baseName))
	nameDoc = genutils.CombineClojureName(ns, fmt.Sprintf(pattern, baseNameDoc))
	goApiString = genutils.CombineClojureName(ns, fmt.Sprintf(pattern, goApiString))

	//	fmt.Printf("jtypes.go/typeNameForExpr: %s (`%s' %s %s) %+v => `%s' %T at:%s\n", name, pattern, ns, baseName, e, pattern, ue, WhereAt(e.Pos()))

	return
}

func Define(ts *TypeSpec, varExpr Expr) *Info {

	ns := ClojureNamespaceForPos(Fset.Position(ts.Name.NamePos))

	pattern, _ := patternForExpr(varExpr)

	name := genutils.CombineClojureName(ns, fmt.Sprintf(pattern, ts.Name.Name))

	jti := &Info{
		FullName:          name,
		FullNameDoc:       name,
		Who:               "TypeDefine",
		Pattern:           pattern,
		Namespace:         ns,
		BaseName:          ts.Name.Name,
		BaseNameDoc:       ts.Name.Name,
		ArgExtractFunc:    "Object",
		ArgClojureArgType: name,
		AsClojureObject:   "GoObjectIfNeeded(%s%s)",
		GoApiString:       name,
	}

	jti.register()

	return jti

}

func InfoForGoName(fullName string) *Info {
	return goTypeMap[fullName]
}

func InfoForExpr(e Expr) *Info {
	if info, ok := typesByExpr[e]; ok {
		return info
	}

	pattern, ns, baseName, baseNameDoc, fullName, fullNameDoc, goApiString, info := namingForExpr(e)

	if info != nil {
		// Already found info on builtin Go type, so just return that.
		typesByExpr[e] = info
		return info
	}

	if inf, found := typesByFullname[fullName]; found {
		typesByExpr[e] = inf
		return inf
	}

	convertFromClojure, convertFromMap := ConversionsFn(e)

	info = &Info{
		Expr:               e,
		FullName:           fullName,
		FullNameDoc:        fullNameDoc,
		Who:                fmt.Sprintf("TypeForExpr %T", e),
		Pattern:            pattern,
		Namespace:          ns,
		BaseName:           baseName,
		BaseNameDoc:        baseNameDoc,
		ArgClojureArgType:  fullName,
		ConvertFromClojure: convertFromClojure,
		ConvertFromMap:     convertFromMap,
		GoApiString:        goApiString,
	}

	typesByExpr[e] = info
	typesByFullname[fullName] = info

	return info
}

func (ti *Info) NameDoc(e Expr) string {
	if ti.Pattern == "" || ti.Namespace == "" {
		return ti.FullNameDoc
	}
	if e != nil && ClojureNamespaceForExpr(e) != ti.Namespace {
		//		fmt.Printf("jtypes.NameDoc(%+v at %s) => %s (in ns=%s) per %s\n", e, WhereAt(e.Pos()), ti.FullName, ti.Namespace, ClojureNamespaceForExpr(e))
		return ti.FullNameDoc
	}
	res := fmt.Sprintf(ti.Pattern, ti.BaseNameDoc)
	//	fmt.Printf("jtypes.NameDoc(%+v at %s) => just %s (in ns=%s) per %s\n", e, WhereAt(e.Pos()), res, ti.Namespace, ClojureNamespaceForExpr(e))
	return res
}

func (ti *Info) register() {
	if _, found := typesByFullname[ti.FullName]; !found {
		typesByFullname[ti.FullName] = ti
	}
}

var Nil = &Info{}

var Error = &Info{
	FullName:             "Error",
	FullNameDoc:          "Error",
	BaseName:             "Error",
	BaseNameDoc:          "Error",
	ArgClojureType:       "Error",
	ArgFromClojureObject: "",
	ArgExtractFunc:       "Error",
	ArgClojureArgType:    "Error",
	ConvertFromMap:       `FieldAs_error(%s, %s)`,
	AsClojureObject:      "Error(%s%s)",
	ConvertFromClojure:   "ObjectAs_error(%s, %s)",
	PromoteType:          "%s",
	GoApiString:          "error",
}

var Boolean = &Info{
	FullName:             "Boolean",
	FullNameDoc:          "Boolean",
	BaseName:             "Boolean",
	BaseNameDoc:          "Boolean",
	ArgClojureType:       "Boolean",
	ArgFromClojureObject: ".B",
	ArgExtractFunc:       "Boolean",
	ArgClojureArgType:    "Boolean",
	ConvertFromMap:       "FieldAs_bool(%s, %s)",
	AsClojureObject:      "Boolean(%s%s)",
	ConvertFromClojure:   "ObjectAs_bool(%s, %s)",
	PromoteType:          "%s",
	GoApiString:          "bool",
}

var Byte = &Info{
	FullName:             "Byte",
	FullNameDoc:          "Byte",
	BaseName:             "Byte",
	BaseNameDoc:          "Byte",
	ArgClojureType:       "Int",
	ArgFromClojureObject: ".Int().I",
	ArgExtractFunc:       "uint8",
	ArgClojureArgType:    "Int",
	ConvertFromMap:       `FieldAs_uint8(%s, %s)`,
	AsClojureObject:      "Int(int(%s)%s)",
	ConvertFromClojure:   "ObjectAs_uint8(%s, %s)",
	PromoteType:          "int(%s)",
	GoApiString:          "uint8",
	AliasOf:              &UInt8,
}

var Rune = &Info{
	FullName:             "Char",
	FullNameDoc:          "Char",
	BaseName:             "Char",
	BaseNameDoc:          "Char",
	ArgClojureType:       "Char",
	ArgFromClojureObject: ".Ch",
	ArgExtractFunc:       "Char",
	ArgClojureArgType:    "Char",
	ConvertFromMap:       `FieldAs_rune(%s, %s)`,
	AsClojureObject:      "Char(%s%s)",
	ConvertFromClojure:   "ObjectAs_rune(%s, %s)",
	PromoteType:          "%s",
	GoApiString:          "rune",
	AliasOf:              &Int32,
}

var String = &Info{
	FullName:             "String",
	FullNameDoc:          "String",
	BaseName:             "String",
	BaseNameDoc:          "String",
	ArgClojureType:       "String",
	ArgFromClojureObject: ".S",
	ArgExtractFunc:       "String",
	ArgClojureArgType:    "String",
	ConvertFromMap:       `FieldAs_string(%s, %s)`,
	AsClojureObject:      "String(%s%s)",
	ConvertFromClojure:   "ObjectAs_string(%s, %s)",
	PromoteType:          "%s",
	GoApiString:          "string",
}

var Int = &Info{
	FullName:             "Int",
	FullNameDoc:          "Int",
	BaseName:             "Int",
	BaseNameDoc:          "Int",
	ArgClojureType:       "Number",
	ArgFromClojureObject: ".Int().I",
	ArgExtractFunc:       "Int",
	ArgClojureArgType:    "Int",
	ConvertFromMap:       `FieldAs_int(%s, %s)`,
	AsClojureObject:      "Int(%s%s)",
	ConvertFromClojure:   "ObjectAs_int(%s, %s)",
	PromoteType:          "%s",
	GoApiString:          "int",
}

var Int8 = &Info{
	FullName:             "Int8",
	FullNameDoc:          "Int",
	BaseName:             "Int8",
	BaseNameDoc:          "Int",
	ArgClojureType:       "Int",
	ArgFromClojureObject: ".Int().I",
	ArgExtractFunc:       "int8",
	ArgClojureArgType:    "Int",
	ConvertFromMap:       `FieldAs_int8(%s, %s)`,
	AsClojureObject:      "Int(int(%s)%s)",
	ConvertFromClojure:   "ObjectAs_int8(%s, %s)",
	PromoteType:          "int(%s)",
	GoApiString:          "int8",
}

var Int16 = &Info{
	FullName:             "Int16",
	FullNameDoc:          "Int",
	BaseName:             "Int16",
	BaseNameDoc:          "Int",
	ArgClojureType:       "Number",
	ArgFromClojureObject: ".Int().I",
	ArgExtractFunc:       "int16",
	ArgClojureArgType:    "Int",
	ConvertFromMap:       `FieldAs_int16(%s, %s)`,
	AsClojureObject:      "Int(int(%s)%s)",
	ConvertFromClojure:   "ObjectAs_int16(%s, %s)",
	PromoteType:          "int(%s)",
	GoApiString:          "int16",
}

var Int32 = &Info{
	FullName:             "Int32",
	FullNameDoc:          "Int",
	BaseName:             "Int32",
	BaseNameDoc:          "Int",
	ArgClojureType:       "Number",
	ArgFromClojureObject: ".Int().I",
	ArgExtractFunc:       "int32",
	ArgClojureArgType:    "Int",
	ConvertFromMap:       `FieldAs_int32(%s, %s)`,
	AsClojureObject:      "Int(int(%s)%s)",
	ConvertFromClojure:   "ObjectAs_int32(%s, %s)",
	PromoteType:          "int(%s)",
	GoApiString:          "int32",
}

var Int64 = &Info{
	FullName:             "Int64",
	FullNameDoc:          "BigInt",
	BaseName:             "Int64",
	BaseNameDoc:          "BigInt",
	ArgClojureType:       "Number",
	ArgFromClojureObject: ".BigInt().Int64()",
	ArgExtractFunc:       "int64",
	ArgClojureArgType:    "BigInt",
	ConvertFromMap:       `FieldAs_int64(%s, %s)`,
	AsClojureObject:      "Number(%s%s)",
	ConvertFromClojure:   "ObjectAs_int64(%s, %s)",
	PromoteType:          "int64(%s)",
	GoApiString:          "int64",
}

var UInt = &Info{
	FullName:             "Uint",
	FullNameDoc:          "Number",
	BaseName:             "Uint",
	BaseNameDoc:          "Number",
	ArgClojureType:       "Number",
	ArgFromClojureObject: ".Int().I",
	ArgExtractFunc:       "uint",
	ArgClojureArgType:    "Number",
	ConvertFromMap:       `FieldAs_uint(%s, %s)`,
	AsClojureObject:      "Number(%s%s)",
	ConvertFromClojure:   "ObjectAs_uint(%s, %s)",
	PromoteType:          "uint64(%s)",
	GoApiString:          "uint",
}

var UInt8 = &Info{
	FullName:             "Uint8",
	FullNameDoc:          "Int",
	BaseName:             "Uint8",
	BaseNameDoc:          "Int",
	ArgClojureType:       "Int",
	ArgFromClojureObject: ".Int().I",
	ArgExtractFunc:       "uint8",
	ArgClojureArgType:    "Int",
	ConvertFromMap:       `FieldAs_uint8(%s, %s)`,
	AsClojureObject:      "Int(int(%s)%s)",
	ConvertFromClojure:   "ObjectAs_uint8(%s, %s)",
	PromoteType:          "int(%s)",
	GoApiString:          "uint8",
}

var UInt16 = &Info{
	FullName:             "Uint16",
	FullNameDoc:          "Int",
	BaseName:             "Uint16",
	BaseNameDoc:          "Int",
	ArgClojureType:       "Number",
	ArgFromClojureObject: ".Int().I",
	ArgExtractFunc:       "uint16",
	ArgClojureArgType:    "Int",
	ConvertFromMap:       `FieldAs_uint16(%s, %s)`,
	AsClojureObject:      "Int(int(%s)%s)",
	ConvertFromClojure:   "ObjectAs_uint16(%s, %s)",
	PromoteType:          "int(%s)",
	GoApiString:          "uint16",
}

var UInt32 = &Info{
	FullName:             "Uint32",
	FullNameDoc:          "Number",
	BaseName:             "Uint32",
	BaseNameDoc:          "Number",
	ArgClojureType:       "Number",
	ArgFromClojureObject: ".Int().I",
	ArgExtractFunc:       "uint32",
	ArgClojureArgType:    "Number",
	ConvertFromMap:       `FieldAs_uint32(%s, %s)`,
	AsClojureObject:      "Number(%s%s)",
	ConvertFromClojure:   "ObjectAs_uint32(%s, %s)",
	PromoteType:          "int64(%s)",
	GoApiString:          "uint32",
}

var UInt64 = &Info{
	FullName:             "Uint64",
	FullNameDoc:          "Number",
	BaseName:             "Uint64",
	BaseNameDoc:          "Number",
	ArgClojureType:       "Number",
	ArgFromClojureObject: ".BigInt().Uint64()",
	ArgExtractFunc:       "uint64",
	ArgClojureArgType:    "Number",
	ConvertFromMap:       `FieldAs_uint64(%s, %s)`,
	AsClojureObject:      "Number(%s%s)",
	ConvertFromClojure:   "ObjectAs_uint64(%s, %s)",
	PromoteType:          "uint64(%s)",
	GoApiString:          "uint64",
}

var UIntPtr = &Info{
	FullName:             "UintPtr",
	FullNameDoc:          "Number",
	BaseName:             "UintPtr",
	BaseNameDoc:          "Number",
	ArgClojureType:       "Number",
	ArgFromClojureObject: ".BigInt().Uint64()",
	ArgExtractFunc:       "uintptr",
	ArgClojureArgType:    "Number",
	ConvertFromMap:       `FieldAs_uintptr(%s, %s)`,
	AsClojureObject:      "Number(%s%s)",
	ConvertFromClojure:   "ObjectAs_uintptr(%s, %s)",
	PromoteType:          "uint64(%s)",
	GoApiString:          "uintptr",
}

var Float32 = &Info{
	FullName:             "Float32",
	FullNameDoc:          "Double",
	BaseName:             "Float32",
	BaseNameDoc:          "Double",
	ArgClojureType:       "",
	ArgFromClojureObject: "",
	ArgExtractFunc:       "float32",
	ArgClojureArgType:    "Double",
	ConvertFromMap:       `FieldAs_float32(%s, %s)`,
	AsClojureObject:      "Double(float64(%s)%s)",
	ConvertFromClojure:   "ObjectAs_float32(%s, %s)",
	PromoteType:          "float64(%s)",
	GoApiString:          "float32",
}

var Float64 = &Info{
	FullName:             "Double",
	FullNameDoc:          "Double",
	BaseName:             "Double",
	BaseNameDoc:          "Double",
	ArgClojureType:       "Double",
	ArgFromClojureObject: "",
	ArgExtractFunc:       "Double",
	ArgClojureArgType:    "Double",
	ConvertFromMap:       `FieldAs_float64(%s, %s)`,
	AsClojureObject:      "Double(%s%s)",
	ConvertFromClojure:   "ObjectAs_float64(%s, %s)",
	PromoteType:          "%s",
	GoApiString:          "float64",
}

var Complex128 = &Info{
	FullName:             "ABEND007(find these)",
	FullNameDoc:          "ABEND007(find these)",
	BaseName:             "ABEND007(find these)",
	BaseNameDoc:          "ABEND007(find these)",
	ArgClojureType:       "",
	ArgFromClojureObject: "",
	ArgExtractFunc:       "complex128",
	ArgClojureArgType:    "ABEND007(find these)",
	ConvertFromMap:       "FieldAs_complex128(%s, %s)",
	AsClojureObject:      "", // TODO: support complex128 in Clojure, even if via just [real imag]
	ConvertFromClojure:   "ObjectAs_complex128(%s, %s)",
	GoApiString:          "complex128",
}

var goTypeMap = map[string]*Info{
	"nil":        Nil,
	"error":      Error,
	"bool":       Boolean,
	"byte":       Byte,
	"rune":       Rune,
	"string":     String,
	"int":        Int,
	"int8":       Int8,
	"int16":      Int16,
	"int32":      Int32,
	"int64":      Int64,
	"uint":       UInt,
	"uint8":      UInt8,
	"uint16":     UInt16,
	"uint32":     UInt32,
	"uint64":     UInt64,
	"uintptr":    UIntPtr,
	"float32":    Float32,
	"float64":    Float64,
	"complex128": Complex128,
}

var ConversionsFn func(e Expr) (fromClojure, fromMap string)
