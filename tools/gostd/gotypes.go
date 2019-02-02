package main

import (
	"fmt"
	. "go/ast"
	gotypes "go/types"
)

type goTypeInfo struct {
	typeName                  string          // empty (not a declared type) or the basic type name ("foo" for "x/y.foo")
	fullName                  string          // empty ("struct {...}" etc.), typeName (built-in), path/to/pkg.typeName, or ABEND if unsupported
	underlyingType            *Expr           // nil if not a declared type
	argClojureType            string          // Can convert this type to a Go function arg with my type
	argFromClojureObject      string          // Append this to Clojure object to extract value of my type
	argClojureArgType         string          // Clojure argument type for a Go function arg with my type
	argExtractFunc            string          // Call Extract<this>() for arg with my type
	convertFromClojure        string          // Pattern to convert a (scalar) %s to this type
	convertFromClojureImports []packageImport // Imports needed to support the above
	uncompleted               bool            // Has this type's info been filled in beyond the registration step?
	custom                    bool            // Is this not a builtin Go type?
	private                   bool            // Is this a private type?
	unsupported               bool            // Is this unsupported?
	constructs                bool            // Does the convertion from Clojure actually construct (via &sometype{}), returning ptr?
}

/* These map fullNames to type info. */
var goTypes = map[string]*goTypeInfo{}

func registerType(gf *goFile, fullTypeName string, ts *TypeSpec) {
	if _, found := goTypes[fullTypeName]; found {
		return
	}
	ti := &goTypeInfo{
		typeName:       ts.Name.Name,
		fullName:       fullTypeName,
		underlyingType: &ts.Type,
		private:        isPrivate(ts.Name.Name),
		custom:         true,
		uncompleted:    true,
	}
	goTypes[fullTypeName] = ti
}

func toGoTypeNameInfo(pkgDirUnix, baseName string, e *Expr) *goTypeInfo {
	if ti, found := goTypes[baseName]; found {
		return ti
	}
	fullName := pkgDirUnix + "." + baseName
	if ti, found := goTypes[fullName]; found {
		return ti
	}
	if gotypes.Universe.Lookup(baseName) != nil {
		ti := &goTypeInfo{
			typeName:           baseName,
			fullName:           fmt.Sprintf("ABEND046(gotypes.go: unsupported builtin type %s)", baseName),
			argClojureType:     baseName,
			argClojureArgType:  baseName,
			convertFromClojure: baseName + "(%s)",
			unsupported:        true,
		}
		goTypes[baseName] = ti
		return ti
	}
	panic(fmt.Sprintf("type %s not found at %s", fullName, whereAt((*e).Pos())))
}

func toGoTypeInfo(src *goFile, ts *TypeSpec) *goTypeInfo {
	return toGoExprInfo(src, &ts.Type)
}

func toGoExprInfo(src *goFile, e *Expr) *goTypeInfo {
	typeName := ""
	fullName := ""
	convertFromClojure := ""
	private := false
	var underlyingType *Expr
	unsupported := false
	switch td := (*e).(type) {
	case *Ident:
		ti := toGoTypeNameInfo(src.pkgDirUnix, td.Name, e)
		if ti.uncompleted {
			// Fill in other info now that all types are registered.
			ut := toGoExprInfo(src, ti.underlyingType)
			if ut.unsupported {
				ti.fullName = ut.fullName
				ti.unsupported = true
			}
			if ut.constructs {
				ti.convertFromClojure = "*" + ut.convertFromClojure
			} else {
				ti.convertFromClojure = ut.convertFromClojure
			}
			ti.convertFromClojureImports = ut.convertFromClojureImports
			ti.uncompleted = false
		}
		return ti
	case *ArrayType:
		return goArrayType(src, &td.Len, &td.Elt)
	case *StarExpr:
		return goStarExpr(src, &td.X)
	}
	if typeName == "" || fullName == "" {
		typeName = fmt.Sprintf("%T", *e)
		fullName = fmt.Sprintf("ABEND047(gotypes.go: unsupported type %s)", typeName)
		unsupported = true
	}
	v := &goTypeInfo{
		typeName:           typeName,
		fullName:           fullName,
		underlyingType:     underlyingType,
		private:            private,
		unsupported:        unsupported,
		convertFromClojure: convertFromClojure,
	}
	goTypes[fullName] = v
	return v
}

func toGoExprString(src *goFile, e *Expr) string {
	if e == nil {
		return "-"
	}
	t := toGoExprInfo(src, e)
	if t != nil {
		return t.fullName
	}
	return fmt.Sprintf("%T", e)
}

func lenString(len *Expr) string {
	if len == nil || *len == nil {
		return ""
	}
	l := *len
	switch n := l.(type) {
	case *Ident:
		return n.Name
	case *BasicLit:
		return n.Value
	}
	return fmt.Sprintf("%T", l)
}

func goArrayType(src *goFile, len *Expr, elt *Expr) *goTypeInfo {
	var fullName string
	e := toGoExprInfo(src, elt)
	fullName = "[" + lenString(len) + "]" + e.fullName
	if v, ok := goTypes[fullName]; ok {
		return v
	}
	v := &goTypeInfo{
		typeName:       e.typeName,
		fullName:       fullName,
		underlyingType: elt,
		custom:         true,
		unsupported:    e.unsupported,
		constructs:     e.constructs,
	}
	goTypes[fullName] = v
	return v
}

func ptrTo(expr string) string {
	if expr[0] == '*' {
		return expr[1:]
	}
	return expr
}

func goStarExpr(src *goFile, x *Expr) *goTypeInfo {
	e := toGoExprInfo(src, x)
	fullName := "*" + e.fullName
	if v, ok := goTypes[fullName]; ok {
		return v
	}
	convertFromClojure := ""
	if e.convertFromClojure != "" {
		if e.constructs {
			convertFromClojure = e.convertFromClojure
		} else {
			convertFromClojure = "&" + e.convertFromClojure
		}
	}
	v := &goTypeInfo{
		typeName:           e.typeName,
		fullName:           fullName,
		underlyingType:     x,
		convertFromClojure: convertFromClojure,
		custom:             true,
		private:            e.private,
		unsupported:        e.unsupported,
	}
	goTypes[fullName] = v
	return v
}

func init() {
	goTypes["bool"] = &goTypeInfo{
		typeName:             "bool",
		fullName:             "bool",
		argClojureType:       "Boolean",
		argFromClojureObject: ".Boolean().B",
		argClojureArgType:    "Boolean",
		argExtractFunc:       "Boolean",
		convertFromClojure:   "ToBool(%s)",
	}
	goTypes["string"] = &goTypeInfo{
		typeName:             "string",
		fullName:             "string",
		argClojureType:       "String",
		argFromClojureObject: ".S",
		argClojureArgType:    "String",
		argExtractFunc:       "String",
		convertFromClojure:   `AssertString(%s, "").S`,
	}
	goTypes["rune"] = &goTypeInfo{
		typeName:             "rune",
		fullName:             "rune",
		argClojureType:       "Char",
		argFromClojureObject: ".Ch",
		argClojureArgType:    "Char",
		argExtractFunc:       "Char",
		convertFromClojure:   `AssertChar(%s, "").Ch`,
	}
	goTypes["byte"] = &goTypeInfo{
		typeName:             "byte",
		fullName:             "byte",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "Byte",
		convertFromClojure:   `byte(AssertInt(%s, "").I)`,
	}
	goTypes["int"] = &goTypeInfo{
		typeName:             "int",
		fullName:             "int",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "Int",
		convertFromClojure:   `AssertInt(%s, "").I`,
	}
	goTypes["uint"] = &goTypeInfo{
		typeName:             "uint",
		fullName:             "uint",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Number",
		argExtractFunc:       "UInt",
		convertFromClojure:   `uint(AssertInt(%s, "").I)`,
	}
	goTypes["int8"] = &goTypeInfo{
		typeName:             "int8",
		fullName:             "int8",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "Byte",
		convertFromClojure:   `int8(AssertInt(%s, "").I)`,
	}
	goTypes["uint8"] = &goTypeInfo{
		typeName:             "uint8",
		fullName:             "uint8",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "UInt8",
		convertFromClojure:   `uint8(AssertInt(%s, "").I)`,
	}
	goTypes["int16"] = &goTypeInfo{
		typeName:             "int16",
		fullName:             "int16",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "Int16",
		convertFromClojure:   `int16(AssertInt(%s, "").I)`,
	}
	goTypes["uint16"] = &goTypeInfo{
		typeName:             "uint16",
		fullName:             "uint16",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "UInt16",
		convertFromClojure:   `uint16(AssertInt(%s, "").I)`,
	}
	goTypes["int32"] = &goTypeInfo{
		typeName:             "int32",
		fullName:             "int32",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Int",
		argExtractFunc:       "Int32",
		convertFromClojure:   `int32(AssertInt(%s, "").I)`,
	}
	goTypes["uint32"] = &goTypeInfo{
		typeName:             "uint32",
		fullName:             "uint32",
		argClojureType:       "Number",
		argFromClojureObject: ".Int().I",
		argClojureArgType:    "Number",
		argExtractFunc:       "UInt32",
		convertFromClojure:   `uint32(AssertNumber(%s, "").BigInt().Uint64())`,
	}
	goTypes["int64"] = &goTypeInfo{
		typeName:             "int64",
		fullName:             "int64",
		argClojureType:       "Number",
		argFromClojureObject: ".BigInt().Int64()",
		argClojureArgType:    "Number",
		argExtractFunc:       "Int64",
		convertFromClojure:   `AssertNumber(%s, "").BigInt().Int64()`,
	}
	goTypes["uint64"] = &goTypeInfo{
		typeName:             "uint64",
		fullName:             "uint64",
		argClojureType:       "Number",
		argFromClojureObject: ".BigInt().Uint64()",
		argClojureArgType:    "Number",
		argExtractFunc:       "UInt64",
		convertFromClojure:   `AssertNumber(%s, "").BigInt().Uint64()`,
	}
	goTypes["uintptr"] = &goTypeInfo{
		typeName:             "uintptr",
		fullName:             "uintptr",
		argClojureType:       "Number",
		argFromClojureObject: ".BigInt().Uint64()",
		argClojureArgType:    "Number",
		argExtractFunc:       "UIntPtr",
		convertFromClojure:   `uintptr(AssertNumber(%s, "").BigInt().Uint64())`,
	}
	goTypes["float32"] = &goTypeInfo{
		typeName:             "float32",
		fullName:             "float32",
		argClojureType:       "Double",
		argFromClojureObject: "",
		argClojureArgType:    "Double",
		argExtractFunc:       "ABEND007(find these)",
		convertFromClojure:   `float32(AssertDouble(%s, "").D)`,
	}
	goTypes["float64"] = &goTypeInfo{
		typeName:             "float64",
		fullName:             "float64",
		argClojureType:       "Double",
		argFromClojureObject: "",
		argClojureArgType:    "Double",
		argExtractFunc:       "ABEND007(find these)",
		convertFromClojure:   `float64(AssertDouble(%s, "").D)`,
	}
	goTypes["complex64"] = &goTypeInfo{
		typeName:             "complex64",
		fullName:             "complex64",
		argClojureType:       "",
		argFromClojureObject: "",
		argClojureArgType:    "",
		argExtractFunc:       "ABEND007(find these)",
		convertFromClojure:   "", // TODO: support this in Joker, even if via just [real imag]
	}
	goTypes["complex128"] = &goTypeInfo{
		typeName:             "complex128",
		fullName:             "complex128",
		argClojureType:       "",
		argFromClojureObject: "",
		argClojureArgType:    "",
		argExtractFunc:       "ABEND007(find these)",
		convertFromClojure:   "", // TODO: support this in Joker, even if via just [real imag]
	}
	goTypes["error"] = &goTypeInfo{
		typeName:                  "error",
		fullName:                  "error",
		argClojureType:            "Error",
		argFromClojureObject:      "",
		argClojureArgType:         "String",
		argExtractFunc:            "",
		convertFromClojure:        `_errors.New(AssertString(%s, "").S)`,
		convertFromClojureImports: []packageImport{{"_errors", "errors"}},
	}
}
