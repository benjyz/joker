package main

import (
	"fmt"
	. "github.com/candid82/joker/tools/gostd/gowalk"
	"github.com/candid82/joker/tools/gostd/imports"
	. "github.com/candid82/joker/tools/gostd/utils"
	. "go/ast"
	gotypes "go/types"
)

func registerType(gf *GoFile, fullGoTypeName string, ts *TypeSpec) *GoTypeInfo {
	if ti, found := GoTypes[fullGoTypeName]; found {
		return ti
	}
	ti := &GoTypeInfo{
		LocalName:         ts.Name.Name,
		FullGoName:        fullGoTypeName,
		SourceFile:        gf,
		UnderlyingType:    &ts.Type,
		ArgClojureArgType: FullTypeNameAsClojure(gf.NsRoot, fullGoTypeName),
		Private:           IsPrivate(ts.Name.Name),
		Custom:            true,
		Uncompleted:       true,
		ConvertToClojure:  "GoObject(%s%s)",
		ArgExtractFunc:    "Object",
	}
	GoTypes[fullGoTypeName] = ti
	return ti
}

func toGoTypeNameInfo(pkgDirUnix, baseName string, e *Expr) *GoTypeInfo {
	if ti, found := GoTypes[baseName]; found {
		return ti
	}
	fullGoName := pkgDirUnix + "." + baseName
	if ti, found := GoTypes[fullGoName]; found {
		return ti
	}
	if gotypes.Universe.Lookup(baseName) != nil {
		ti := &GoTypeInfo{
			LocalName:          baseName,
			FullGoName:         fmt.Sprintf("ABEND046(gotypes.go: unsupported builtin type %s for %s)", baseName, pkgDirUnix),
			ArgClojureType:     baseName,
			ArgClojureArgType:  baseName,
			ConvertFromClojure: baseName + "(%s)",
			ConvertToClojure:   "GoObject(%s%s)",
			Unsupported:        true,
		}
		GoTypes[baseName] = ti
		return ti
	}
	panic(fmt.Sprintf("type %s not found at %s", fullGoName, WhereAt((*e).Pos())))
}

func toGoTypeInfo(src *GoFile, ts *TypeSpec) *GoTypeInfo {
	return toGoExprInfo(src, &ts.Type)
}

func toGoExprInfo(src *GoFile, e *Expr) *GoTypeInfo {
	localName := ""
	fullGoName := ""
	convertFromClojure := ""
	private := false
	var underlyingType *Expr
	unsupported := false
	switch td := (*e).(type) {
	case *Ident:
		ti := toGoTypeNameInfo(src.PkgDirUnix, td.Name, e)
		if ti.Uncompleted {
			// Fill in other info now that all types are registered.
			ut := toGoExprInfo(src, ti.UnderlyingType)
			/*			if ut.unsupported {
						ti.fullGoName = ut.fullGoName
						ti.unsupported = true
					}*/
			if ut.ConvertFromClojure != "" {
				if ut.Constructs {
					ti.ConvertFromClojure = "*" + ut.ConvertFromClojure
				} else {
					ti.ConvertFromClojure = fmt.Sprintf("_%s.%s(%s)", ti.SourceFile.PkgBaseName, ti.LocalName, ut.ConvertFromClojure)
				}
			}
			ti.ConvertFromClojureImports = ut.ConvertFromClojureImports
			ti.Uncompleted = false
		}
		return ti
	case *ArrayType:
		return goArrayType(src, &td.Len, &td.Elt)
	case *StarExpr:
		return goStarExpr(src, &td.X)
	}
	if localName == "" || fullGoName == "" {
		localName = fmt.Sprintf("%T", *e)
		fullGoName = fmt.Sprintf("ABEND047(gotypes.go: unsupported type %s)", localName)
		unsupported = true
	}
	v := &GoTypeInfo{
		LocalName:          localName,
		FullGoName:         fullGoName,
		UnderlyingType:     underlyingType,
		Private:            private,
		Unsupported:        unsupported,
		ConvertFromClojure: convertFromClojure,
		ConvertToClojure:   "GoObject(%s%s)",
	}
	GoTypes[fullGoName] = v
	return v
}

func toGoExprString(src *GoFile, e *Expr) string {
	if e == nil {
		return "-"
	}
	t := toGoExprInfo(src, e)
	if t != nil {
		return t.FullGoName
	}
	return fmt.Sprintf("%T", e)
}

func toGoExprTypeName(src *GoFile, e *Expr) string {
	if e == nil {
		return "-"
	}
	t := toGoExprInfo(src, e)
	if t != nil {
		return t.LocalName
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

func goArrayType(src *GoFile, len *Expr, elt *Expr) *GoTypeInfo {
	var fullGoName string
	e := toGoExprInfo(src, elt)
	fullGoName = "[" + lenString(len) + "]" + e.FullGoName
	if v, ok := GoTypes[fullGoName]; ok {
		return v
	}
	v := &GoTypeInfo{
		LocalName:        e.LocalName,
		FullGoName:       fullGoName,
		UnderlyingType:   elt,
		Custom:           true,
		Unsupported:      e.Unsupported,
		Constructs:       e.Constructs,
		ConvertToClojure: "GoObject(%s%s)",
	}
	GoTypes[fullGoName] = v
	return v
}

func ptrTo(expr string) string {
	if expr[0] == '*' {
		return expr[1:]
	}
	return expr
}

func goStarExpr(src *GoFile, x *Expr) *GoTypeInfo {
	e := toGoExprInfo(src, x)
	fullGoName := "*" + e.FullGoName
	if v, ok := GoTypes[fullGoName]; ok {
		return v
	}
	convertFromClojure := ""
	if e.ConvertFromClojure != "" {
		if e.Constructs {
			convertFromClojure = e.ConvertFromClojure
		} else if e.ArgClojureArgType == e.ArgExtractFunc {
			/* Not a conversion, so can take address of the Clojure object's internals. */
			convertFromClojure = "&" + e.ConvertFromClojure
		}
	}
	v := &GoTypeInfo{
		LocalName:          e.LocalName,
		FullGoName:         fullGoName,
		UnderlyingType:     x,
		ConvertFromClojure: convertFromClojure,
		Custom:             true,
		Private:            e.Private,
		Unsupported:        e.Unsupported,
		ConvertToClojure:   "GoObject(%s%s)",
	}
	GoTypes[fullGoName] = v
	return v
}

func init() {
	GoTypes["bool"] = &GoTypeInfo{
		LocalName:            "bool",
		FullGoName:           "bool",
		ArgClojureType:       "Boolean",
		ArgFromClojureObject: ".B",
		ArgClojureArgType:    "Boolean",
		ArgExtractFunc:       "Boolean",
		ConvertFromClojure:   "ToBool(%s)",
		ConvertToClojure:     "Boolean(%s%s)",
		PromoteType:          "%s",
	}
	GoTypes["string"] = &GoTypeInfo{
		LocalName:            "string",
		FullGoName:           "string",
		ArgClojureType:       "String",
		ArgFromClojureObject: ".S",
		ArgClojureArgType:    "String",
		ArgExtractFunc:       "String",
		ConvertFromClojure:   `AssertString(%s, "").S`,
		ConvertToClojure:     "String(%s%s)",
		PromoteType:          "%s",
	}
	GoTypes["rune"] = &GoTypeInfo{
		LocalName:            "rune",
		FullGoName:           "rune",
		ArgClojureType:       "Char",
		ArgFromClojureObject: ".Ch",
		ArgClojureArgType:    "Char",
		ArgExtractFunc:       "Char",
		ConvertFromClojure:   `AssertChar(%s, "").Ch`,
		ConvertToClojure:     "Char(%s%s)",
		PromoteType:          "%s",
	}
	GoTypes["byte"] = &GoTypeInfo{
		LocalName:            "byte",
		FullGoName:           "byte",
		ArgClojureType:       "Number",
		ArgFromClojureObject: ".Int().I",
		ArgClojureArgType:    "Int",
		ArgExtractFunc:       "Byte",
		ConvertFromClojure:   `byte(AssertInt(%s, "").I)`,
		ConvertToClojure:     "Int(int(%s)%s)",
		PromoteType:          "int(%s)",
	}
	GoTypes["int"] = &GoTypeInfo{
		LocalName:            "int",
		FullGoName:           "int",
		ArgClojureType:       "Number",
		ArgFromClojureObject: ".Int().I",
		ArgClojureArgType:    "Int",
		ArgExtractFunc:       "Int",
		ConvertFromClojure:   `AssertInt(%s, "").I`,
		ConvertToClojure:     "Int(%s%s)",
		PromoteType:          "%s",
	}
	GoTypes["uint"] = &GoTypeInfo{
		LocalName:            "uint",
		FullGoName:           "uint",
		ArgClojureType:       "Number",
		ArgFromClojureObject: ".Int().I",
		ArgClojureArgType:    "Number",
		ArgExtractFunc:       "UInt",
		ConvertFromClojure:   `uint(AssertInt(%s, "").I)`,
		ConvertToClojure:     "BigIntU(uint64(%s)%s)",
		PromoteType:          "uint64(%s)",
	}
	GoTypes["int8"] = &GoTypeInfo{
		LocalName:            "int8",
		FullGoName:           "int8",
		ArgClojureType:       "Int",
		ArgFromClojureObject: ".Int().I",
		ArgClojureArgType:    "Int",
		ArgExtractFunc:       "Int8",
		ConvertFromClojure:   `int8(AssertInt(%s, "").I)`,
		ConvertToClojure:     "Int(int(%s)%s)",
		PromoteType:          "int(%s)",
	}
	GoTypes["uint8"] = &GoTypeInfo{
		LocalName:            "uint8",
		FullGoName:           "uint8",
		ArgClojureType:       "Int",
		ArgFromClojureObject: ".Int().I",
		ArgClojureArgType:    "Int",
		ArgExtractFunc:       "UInt8",
		ConvertFromClojure:   `uint8(AssertInt(%s, "").I)`,
		ConvertToClojure:     "Int(int(%s)%s)",
		PromoteType:          "int(%s)",
	}
	GoTypes["int16"] = &GoTypeInfo{
		LocalName:            "int16",
		FullGoName:           "int16",
		ArgClojureType:       "Number",
		ArgFromClojureObject: ".Int().I",
		ArgClojureArgType:    "Int",
		ArgExtractFunc:       "Int16",
		ConvertFromClojure:   `int16(AssertInt(%s, "").I)`,
		ConvertToClojure:     "Int(int(%s)%s)",
		PromoteType:          "int(%s)",
	}
	GoTypes["uint16"] = &GoTypeInfo{
		LocalName:            "uint16",
		FullGoName:           "uint16",
		ArgClojureType:       "Number",
		ArgFromClojureObject: ".Int().I",
		ArgClojureArgType:    "Int",
		ArgExtractFunc:       "UInt16",
		ConvertFromClojure:   `uint16(AssertInt(%s, "").I)`,
		ConvertToClojure:     "Int(int(%s)%s)",
		PromoteType:          "int(%s)",
	}
	GoTypes["int32"] = &GoTypeInfo{
		LocalName:            "int32",
		FullGoName:           "int32",
		ArgClojureType:       "Number",
		ArgFromClojureObject: ".Int().I",
		ArgClojureArgType:    "Int",
		ArgExtractFunc:       "Int32",
		ConvertFromClojure:   `int32(AssertInt(%s, "").I)`,
		ConvertToClojure:     "Int(int(%s)%s)",
		PromoteType:          "int(%s)",
	}
	GoTypes["uint32"] = &GoTypeInfo{
		LocalName:            "uint32",
		FullGoName:           "uint32",
		ArgClojureType:       "Number",
		ArgFromClojureObject: ".Int().I",
		ArgClojureArgType:    "Number",
		ArgExtractFunc:       "UInt32",
		ConvertFromClojure:   `uint32(AssertNumber(%s, "").BigInt().Uint64())`,
		ConvertToClojure:     "BigIntU(uint64(%s)%s)",
		PromoteType:          "int64(%s)",
	}
	GoTypes["int64"] = &GoTypeInfo{
		LocalName:            "int64",
		FullGoName:           "int64",
		ArgClojureType:       "Number",
		ArgFromClojureObject: ".BigInt().Int64()",
		ArgClojureArgType:    "Number",
		ArgExtractFunc:       "Int64",
		ConvertFromClojure:   `AssertNumber(%s, "").BigInt().Int64()`,
		ConvertToClojure:     "BigInt(%s%s)",
		PromoteType:          "int64(%s)", // constants are not auto-promoted, so promote them explicitly for MakeNumber()
	}
	GoTypes["uint64"] = &GoTypeInfo{
		LocalName:            "uint64",
		FullGoName:           "uint64",
		ArgClojureType:       "Number",
		ArgFromClojureObject: ".BigInt().Uint64()",
		ArgClojureArgType:    "Number",
		ArgExtractFunc:       "UInt64",
		ConvertFromClojure:   `AssertNumber(%s, "").BigInt().Uint64()`,
		ConvertToClojure:     "BigIntU(%s%s)",
		PromoteType:          "uint64(%s)", // constants are not auto-promoted, so promote them explicitly for MakeNumber()
	}
	GoTypes["uintptr"] = &GoTypeInfo{
		LocalName:            "uintptr",
		FullGoName:           "uintptr",
		ArgClojureType:       "Number",
		ArgFromClojureObject: ".BigInt().Uint64()",
		ArgClojureArgType:    "Number",
		ArgExtractFunc:       "UIntPtr",
		ConvertFromClojure:   `uintptr(AssertNumber(%s, "").BigInt().Uint64())`,
		PromoteType:          "int64(%s)",
	}
	GoTypes["float32"] = &GoTypeInfo{
		LocalName:            "float32",
		FullGoName:           "float32",
		ArgClojureType:       "Double",
		ArgFromClojureObject: "",
		ArgClojureArgType:    "Double",
		ArgExtractFunc:       "ABEND007(find these)",
		ConvertFromClojure:   `float32(AssertDouble(%s, "").D)`,
		PromoteType:          "double(%s)",
	}
	GoTypes["float64"] = &GoTypeInfo{
		LocalName:            "float64",
		FullGoName:           "float64",
		ArgClojureType:       "Double",
		ArgFromClojureObject: "",
		ArgClojureArgType:    "Double",
		ArgExtractFunc:       "ABEND007(find these)",
		ConvertFromClojure:   `float64(AssertDouble(%s, "").D)`,
		PromoteType:          "%s",
	}
	GoTypes["complex64"] = &GoTypeInfo{
		LocalName:            "complex64",
		FullGoName:           "complex64",
		ArgClojureType:       "",
		ArgFromClojureObject: "",
		ArgClojureArgType:    "",
		ArgExtractFunc:       "ABEND007(find these)",
		ConvertFromClojure:   "", // TODO: support this in Joker, even if via just [real imag]
	}
	GoTypes["complex128"] = &GoTypeInfo{
		LocalName:            "complex128",
		FullGoName:           "complex128",
		ArgClojureType:       "",
		ArgFromClojureObject: "",
		ArgClojureArgType:    "",
		ArgExtractFunc:       "ABEND007(find these)",
		ConvertFromClojure:   "", // TODO: support this in Joker, even if via just [real imag]
	}
	GoTypes["error"] = &GoTypeInfo{
		LocalName:                 "error",
		FullGoName:                "error",
		ArgClojureType:            "Error",
		ArgFromClojureObject:      "",
		ArgClojureArgType:         "Error",
		ArgExtractFunc:            "Error",
		ConvertFromClojure:        `_errors.New(AssertString(%s, "").S)`,
		ConvertFromClojureImports: []imports.Import{{Local: "_errors", LocalRef: "_errors", Full: "errors"}},
		ConvertToClojure:          "Error(%s%s)",
		Nullable:                  true,
	}
}

func init() {
	RegisterType_func = registerType // TODO: Remove this kludge (allowing gowalk to call this fn) when able
}
