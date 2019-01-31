package main

import (
	"fmt"
	. "go/ast"
)

type goTypeInfo struct {
	fullName             string
	argClojureType       string // Can convert this type to a Go function arg with my type
	argFromClojureObject string // Append this to Clojure object to extract value of my type
	argClojureArgType    string // Clojure argument type for a Go function arg with my type
	argExtractFunc       string // Call Extract<this>() for arg with my type
	builtin              bool   // Is this a builtin Go type?
	private              bool   // Is this a private type?
}

var goBuiltinTypes = map[string]*goTypeInfo{}
var goTypes = map[string]*goTypeInfo{}

func toGoTypeInfo(ts *TypeSpec) *goTypeInfo {
	return toGoExprInfo(ts.Type)
}

func toGoExprInfo(e Expr) *goTypeInfo {
	fullName := fmt.Sprintf("<notfound>%T</notfound>", e)
	switch td := e.(type) {
	case *Ident:
		fullName = td.Name
		v := goBuiltinTypes[td.Name]
		if v != nil {
			return v
		}
	case *ArrayType:
		return goArrayType(td.Len, td.Elt)
	case *StarExpr:
		return goStarExpr(td.X)
	}
	v := &goTypeInfo{
		fullName:             fullName,
		argClojureType:       "",
		argFromClojureObject: "",
	}
	goTypes[v.fullName] = v
	return v
}

func toGoExprString(e Expr) string {
	t := toGoExprInfo(e)
	if t != nil {
		return t.fullName
	}
	return fmt.Sprintf("%T", e)
}

func goArrayType(len Expr, elt Expr) *goTypeInfo {
	var fullName string
	ev := toGoExprInfo(elt)
	en := toGoExprString(elt)
	if len == nil {
		fullName = "[]" + en
	} else {
		fullName = "..." + en
	}
	if v, ok := goTypes[fullName]; ok {
		return v
	}
	v := &goTypeInfo{
		fullName:             fullName,
		argClojureType:       ev.argClojureType,
		argFromClojureObject: ev.argFromClojureObject,
	}
	goTypes[fullName] = v
	return v
}

func goStarExpr(x Expr) *goTypeInfo {
	ex := toGoExprInfo(x)
	fullName := "*" + ex.fullName
	if v, ok := goTypes[fullName]; ok {
		return v
	}
	v := &goTypeInfo{
		fullName:             fullName,
		argClojureType:       ex.argClojureType,
		argFromClojureObject: ex.argFromClojureObject,
	}
	goTypes[fullName] = v
	return v
}

func init() {
	goBuiltinTypes["string"] = &goTypeInfo{
		fullName:             "string",
		argClojureType:       "String",
		argFromClojureObject: ".S",
		argClojureArgType:    "String",
		argExtractFunc:       "String",
		builtin:              true,
	}
	goBuiltinTypes["bool"] = &goTypeInfo{
		fullName:             "bool",
		argClojureType:       "Boolean",
		argFromClojureObject: ".Boolean().B",
		argClojureArgType:    "Boolean",
		argExtractFunc:       "Boolean",
		builtin:              true,
	}
	goBuiltinTypes["int"] = &goTypeInfo{
		fullName:             "int",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "Int",
		builtin:              true,
	}
	goBuiltinTypes["byte"] = &goTypeInfo{
		fullName:             "byte",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "Byte",
		builtin:              true,
	}
	goBuiltinTypes["int8"] = &goTypeInfo{
		fullName:             "int8",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "Byte",
		builtin:              true,
	}
	goBuiltinTypes["int16"] = &goTypeInfo{
		fullName:             "int16",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "Int16",
		builtin:              true,
	}
	goBuiltinTypes["uint"] = &goTypeInfo{
		fullName:             "uint",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Number",
		argExtractFunc:       "UInt",
		builtin:              true,
	}
	goBuiltinTypes["uint8"] = &goTypeInfo{
		fullName:             "uint8",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "UInt8",
		builtin:              true,
	}
	goBuiltinTypes["uint16"] = &goTypeInfo{
		fullName:             "uint16",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "UInt16",
		builtin:              true,
	}
	goBuiltinTypes["int32"] = &goTypeInfo{
		fullName:             "int32",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "Int32",
		builtin:              true,
	}
	goBuiltinTypes["uint32"] = &goTypeInfo{
		fullName:             "uint32",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Number",
		argExtractFunc:       "UInt32",
		builtin:              true,
	}
	goBuiltinTypes["int64"] = &goTypeInfo{
		fullName:             "int64",
		argClojureType:       "Number",
		argFromClojureObject: ".BigInt().Int64()",
		argClojureArgType:    "Number",
		argExtractFunc:       "Int64",
		builtin:              true,
	}
	goBuiltinTypes["uint64"] = &goTypeInfo{
		fullName:             "uint64",
		argClojureType:       "Number",
		argFromClojureObject: ".BigInt().Uint64()",
		argClojureArgType:    "Number",
		argExtractFunc:       "UInt64",
		builtin:              true,
	}
	goBuiltinTypes["uintptr"] = &goTypeInfo{
		fullName:             "uintptr",
		argClojureType:       "Number",
		argFromClojureObject: ".BigInt().Uint64()",
		argClojureArgType:    "Number",
		argExtractFunc:       "UIntPtr",
		builtin:              true,
	}
	goBuiltinTypes["float32"] = &goTypeInfo{
		fullName:             "float32",
		argClojureType:       "float32",
		argFromClojureObject: "",
		argClojureArgType:    "Double",
		argExtractFunc:       "ABEND007(find these)",
		builtin:              true,
	}
	goBuiltinTypes["float64"] = &goTypeInfo{
		fullName:             "float64",
		argClojureType:       "float64",
		argFromClojureObject: "",
		argClojureArgType:    "Double",
		argExtractFunc:       "ABEND007(find these)",
		builtin:              true,
	}
	goBuiltinTypes["complex64"] = &goTypeInfo{
		fullName:             "complex64",
		argClojureType:       "complex64",
		argFromClojureObject: "",
		argClojureArgType:    "",
		argExtractFunc:       "ABEND007(find these)",
		builtin:              true,
	}
	goBuiltinTypes["complex128"] = &goTypeInfo{
		fullName:             "complex128",
		argClojureType:       "complex128",
		argFromClojureObject: "",
		argClojureArgType:    "",
		argExtractFunc:       "ABEND007(find these)",
		builtin:              true,
	}
}
