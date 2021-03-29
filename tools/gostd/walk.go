package main

import (
	"bytes"
	"fmt"
	"github.com/candid82/joker/tools/gostd/astutils"
	"github.com/candid82/joker/tools/gostd/genutils"
	"github.com/candid82/joker/tools/gostd/godb"
	"github.com/candid82/joker/tools/gostd/imports"
	"github.com/candid82/joker/tools/gostd/paths"
	. "go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	. "strings"
)

var WalkDump bool

var NumFunctions int
var NumStandalones int
var NumReceivers int
var NumTypes int
var NumConstants int
var NumVariables int
var NumCtableTypes int
var NumGeneratedFunctions int
var NumGeneratedStandalones int
var NumGeneratedReceivers int
var NumGeneratedConstants int
var NumGeneratedVariables int
var NumGeneratedCtors int

type PackageInfo struct {
	DirUnix          string
	BaseName         string
	ImportsNative    *imports.Imports
	ImportsAutoGen   *imports.Imports
	Pkg              *Package
	NonEmpty         bool   // Whether any non-comment code has been generated
	HasGoFiles       bool   // Whether any .go files (would) have been generated
	ClojureNameSpace string // E.g.: "go.std.net", "x.y.z.whatever"
}

/* Map (Unix-style) relative path to package info */
var PackagesInfo = map[string]*PackageInfo{}

/* Sort the packages -- currently appears to not actually be
/* necessary, probably because of how walkDirs() works. */
func SortedPackagesInfo(m map[string]*PackageInfo, f func(k string, i *PackageInfo)) {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f(k, m[k])
	}
}

type FnCodeInfo struct {
	SourceFile *godb.GoFile
	FnCode     string
	FnDecl     *FuncDecl // Empty for standalones and methods; used to get docstring for receivers
	Params     *FieldList
	FnDoc      *CommentGroup
}

type fnCodeMap map[string]*FnCodeInfo

type CodeInfo struct {
	Constants GoConstantsMap
	Variables GoVariablesMap
	Functions fnCodeMap
	Types     TypesMap
	InitTypes map[TypeInfo]struct{}               // types to be initialized
	InitVars  map[TypeInfo]map[string]*FnCodeInfo // func initNative()'s "info_key1 = ... { key2: value, ... }"
}

/* Map relative (Unix-style) package names to maps of function names to code info and strings. */
var ClojureCode = map[string]CodeInfo{}
var ClojureCodeForType = map[TypeInfo]string{}
var GoCode = map[string]CodeInfo{}
var GoCodeForType = map[TypeInfo]string{}

func SortedPackageMap(m map[string]CodeInfo, f func(k string, v CodeInfo)) {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f(k, m[k])
	}
}

func SortedCodeMap(m CodeInfo, f func(k string, v *FnCodeInfo)) {
	var keys []string
	for k, _ := range m.Functions {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f(k, m.Functions[k])
	}
}

func SortedFnCodeInfo(m map[string]*FnCodeInfo, f func(k string, v *FnCodeInfo)) {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f(k, m[k])
	}
}

type FuncInfo struct {
	BaseName       string    // Just the name without receiver-type info
	ReceiverId     string    // Receiver info (only one type supported here and by Golang itself for now)
	Name           string    // Unique name for implementation (has Receiver info as a prefix, then baseName)
	DocName        string    // Everything, for documentation and diagnostics
	Fd             *FuncDecl // nil for methods
	ToM            TypeInfo  // Method operates on this type (nil for standalones and receivers)
	Ft             *FuncType
	Doc            *CommentGroup
	SourceFile     *godb.GoFile
	ImportsNative  *imports.Imports // Add these to package imports if function is generated (no ABENDs)
	ImportsAutoGen *imports.Imports // Add these to package imports if function is generated (no ABENDs)
	Pos            token.Pos
}

func initPackage(rootUnix, pkgDirUnix, nsRoot string, p *Package) {
	if godb.Verbose {
		genutils.AddSortedStdout(fmt.Sprintf("Processing package=%s:\n", pkgDirUnix))
	}

	if _, ok := PackagesInfo[pkgDirUnix]; !ok {
		PackagesInfo[pkgDirUnix] = &PackageInfo{pkgDirUnix, filepath.Base(pkgDirUnix), &imports.Imports{}, &imports.Imports{},
			p, false, false, godb.ClojureNamespaceForDirname(pkgDirUnix)}
		GoCode[pkgDirUnix] = CodeInfo{GoConstantsMap{}, GoVariablesMap{}, fnCodeMap{}, TypesMap{},
			map[TypeInfo]struct{}{}, map[TypeInfo]map[string]*FnCodeInfo{}}
		ClojureCode[pkgDirUnix] = CodeInfo{GoConstantsMap{}, GoVariablesMap{}, fnCodeMap{}, TypesMap{},
			map[TypeInfo]struct{}{}, map[TypeInfo]map[string]*FnCodeInfo{}}
	}
}

/* Go apparently doesn't support/allow 'interface{}' as the value (or
/* key) of a map such that any arbitrary type can be substituted at
/* run time, so there are several of these nearly-identical functions
/* sprinkled through this code. Still get some reuse out of some of
/* them, and it's still easier to maintain these copies than if the
/* body of these were to be included at each call point.... */
func SortedFuncInfoMap(m map[string]*FuncInfo, f func(k string, v *FuncInfo)) {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f(k, m[k])
	}
}

// Add whatever ti needs to be code-generated for fn to fn's list of
// imports; return what is picked as the Go short package name for the
// generated file.
func (fn *FuncInfo) AddToImports(ti TypeInfo) string {
	exprPkgName := ti.GoPackage()
	curPkgName := fn.SourceFile.Package.Dir
	if exprPkgName == "" {
		return ""
	}
	clojureStdNs := "joker.std." + fn.SourceFile.Package.NsRoot
	clojureStdPath := godb.ClojureSourceDir.Join(importStdRoot.String(), goStdPrefix.String()).ToUnix().String()

	native := fn.ImportsNative.AddPackage(exprPkgName, clojureStdNs, clojureStdPath, "", true, fn.Pos)
	if curPkgName.String() == ti.GoPackage() {
		return native
	}
	autoGen := fn.ImportsAutoGen.AddPackage(exprPkgName, clojureStdNs, clojureStdPath, "", true, fn.Pos)
	if native != autoGen {
		panic(fmt.Sprintf("disagreement over '%s': native='%s' autoGen='%s'", exprPkgName, native, autoGen))
	}
	return native
}

func (fn *FuncInfo) AddApiToImports(clType string) string {
	ix := Index(clType, "/")
	if ix < 0 {
		return "" // builtin type (api is in core)
	}

	apiPkgPath := godb.ClojureSourceDir.Join(importStdRoot.String(), ReplaceAll(clType[0:ix], ".", "/")).String()
	clojureStdPath := godb.ClojureSourceDir.Join(importStdRoot.String()).String()
	//	fmt.Fprintf(os.Stderr, "walk.go/AddApiToImports: Compared %s to %s\n", apiPkgPath, fn.SourceFile.Package.ImportMe)
	if apiPkgPath == fn.SourceFile.Package.ImportMe {
		return "" // api is local to function
	}

	clojureStdNs := fn.SourceFile.Package.NsRoot
	native := fn.ImportsNative.AddPackage(apiPkgPath, clojureStdNs, clojureStdPath, "_gostd", true, fn.Pos)

	return native
}

func processTypeRef(t Expr) {
	//	fmt.Printf("%T\n", t)
	if t != nil {
		TypeInfoForExpr(t)
	}
}

func processFieldsForTypes(items []astutils.FieldItem) {
	for _, f := range items {
		processTypeRef(f.Field.Type)
	}
}

func declFuncForTypes(gf *godb.GoFile, pkgDirUnix string, f *File, fd *FuncDecl) {
	if !IsExported(fd.Name.Name) {
		return // Skipping non-exported functions
	}

	processFieldsForTypes(astutils.FlattenFieldList(fd.Recv))
	processFieldsForTypes(astutils.FlattenFieldList(fd.Type.Params))
	processFieldsForTypes(astutils.FlattenFieldList(fd.Type.Results))
}

func processValueSpecsForTypes(gf *godb.GoFile, pkg string, tss []Spec, parentDoc *CommentGroup) {
	for _, spec := range tss {
		ts := spec.(*ValueSpec)
		processTypeRef(ts.Type)
	}
}

// Map qualified function names to info on each.
var QualifiedFunctions = map[string]*FuncInfo{}

func receiverPrefix(src *godb.GoFile, rl []astutils.FieldItem) string {
	res := ""
	for i, r := range rl {
		if i != 0 {
			res += "_"
		}
		switch x := r.Field.Type.(type) {
		case *Ident:
			res += x.Name
		case *ArrayType:
			res += "ArrayOf_" + x.Elt.(*Ident).Name
		case *StarExpr:
			res += "PtrTo_" + x.X.(*Ident).Name
		default:
			panic(fmt.Sprintf("receiverList: unrecognized expr %T in %s", x, src.Name))
		}
	}
	return res + "_"
}

func receiverId(src *godb.GoFile, pkgName string, rl []astutils.FieldItem) string {
	pkg := "{{myGoImport}}."
	res := ""
	for i, r := range rl {
		if i != 0 {
			res += "ABEND422(more than one receiver in list)"
		}
		switch x := r.Field.Type.(type) {
		case *Ident:
			res += pkg + x.Name
		case *ArrayType:
			res += "[]" + pkg + x.Elt.(*Ident).Name
		case *StarExpr:
			res += "*" + pkg + x.X.(*Ident).Name
		default:
			panic(fmt.Sprintf("receiverId: unrecognized expr %T in %s", x, src.Name))
		}
	}
	return res
}

func processFuncDecl(gf *godb.GoFile, pkgDirUnix string, f *File, fd *FuncDecl) {
	if WalkDump {
		fmt.Printf("Func in pkgDirUnix=%s filename=%s fd=%p fd.Doc=%p:\n", pkgDirUnix, godb.FileAt(fd.Pos()), fd, fd.Doc)
		Print(godb.Fset, fd)
	}
	fl := astutils.FlattenFieldList(fd.Recv)
	fnName := receiverPrefix(gf, fl) + fd.Name.Name
	fullName := pkgDirUnix + "." + fnName
	if fullName == "unsafe._Offsetof" {
		return // unsafe.Offsetof is a syntactic operation in Go.
	}
	if v, ok := QualifiedFunctions[fullName]; ok {
		genutils.AddSortedStdout(fmt.Sprintf("NOTE: Already seen function %s in %s, yet again in %s",
			fullName, v.SourceFile.Name, godb.FileAt(fd.Pos())))
	}
	rcvrId := receiverId(gf, gf.Package.BaseName, fl)
	docName := "(" + receiverId(gf, pkgDirUnix, fl) + ")" + fd.Name.Name + "()"
	QualifiedFunctions[fullName] = &FuncInfo{fd.Name.Name, rcvrId, fnName, docName, fd, nil, fd.Type, fd.Doc, gf, &imports.Imports{}, &imports.Imports{}, fd.Pos()}
}

func processTypeDecls(gf *godb.GoFile, pkg string, tss []Spec, parentDoc *CommentGroup) {
	for _, spec := range tss {
		ts := spec.(*TypeSpec)
		RegisterTypeDecl(ts, gf, pkg, parentDoc)
	}
}

func processTypesForTypeDecls(gf *godb.GoFile, pkg string, tss []Spec, parentDoc *CommentGroup) {
	for _, spec := range tss {
		ts := spec.(*TypeSpec)
		RegisterAllSubtypes(ts.Type)
	}
}

type VariableInfo struct {
	Name       *Ident
	SourceFile *godb.GoFile
	Def        string
}

type GoVariablesMap map[string]*VariableInfo

var GoVariables = GoVariablesMap{}

func SortedVariableInfoMap(m map[string]*VariableInfo, f func(k string, v *VariableInfo)) {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f(k, m[k])
	}
}

type ConstantInfo struct {
	Name       *Ident
	SourceFile *godb.GoFile
	Def        string
}

type GoConstantsMap map[string]*ConstantInfo

var GoConstants = GoConstantsMap{}

func SortedConstantInfoMap(m map[string]*ConstantInfo, f func(k string, v *ConstantInfo)) {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f(k, m[k])
	}
}

func fitInt(value string) string {
	_, e := strconv.ParseInt(value, 0, 32)
	if e == nil {
		return "int"
	}
	_, e = strconv.ParseInt(value, 0, 64)
	if e == nil {
		return "int64"
	}
	_, e = strconv.ParseUint(value, 0, 64)
	if e == nil {
		return "uint64"
	}
	return ""
}

func evalConstType(ty *TypeSpec) (typeName string) {
	typeName = ty.Name.Name
	ti := TypeInfoForGoName(typeName)
	if ti == nil {
		// Not a known type; use the underlying type.
		typeName = ty.Type.(*Ident).Name
	}

	return
}

func evalConstExpr(val Expr) (typeName, result string) {
	switch v := val.(type) {
	case *BasicLit:
		result = v.Value
		switch v.Kind {
		case token.STRING:
			typeName = "string"
		case token.INT:
			typeName = fitInt(result)
		case token.FLOAT:
			typeName = "float64"
		case token.CHAR:
			typeName = "rune"
		}
	case *UnaryExpr:
		typeName, result = evalConstExpr(v.X)
		if result == "" && typeName == "" {
			typeName = "int" // TODO: maybe not, but try this for now
		}
		switch v.Op {
		case token.SUB:
			typeName, result = fitInt(result), "-"+result
		default:
		}
	case *BinaryExpr:
		leftType, _ := evalConstExpr(v.X)
		rightType, rightValue := evalConstExpr(v.Y)
		if leftType == rightType {
			typeName = leftType
		} else if leftType == "float64" || rightType == "float64" {
			typeName = "float64" // TODO: probably a good guess for now
		} else if leftType == "int64" || rightType == "int64" {
			typeName = "int64"
		} else if leftType == "rune" || rightType == "rune" {
			typeName = "int"
		}
		if typeName == "int" && v.Op == token.SHL {
			if rightValue == "64" { // TODO: this supports MaxUint64 but is overly specific
				typeName, result = "uint64", strconv.FormatUint(math.MaxUint64, 10)
			} else {
				typeName = "int64"
			}
		} else if typeName == "" {
			typeName = "uint64" // TODO: the outer MaxUint64 definition
		}
	case *ParenExpr:
		typeName, result = evalConstExpr(v.X)
	case *Ident:
		switch v.Name {
		case "iota":
			typeName, result = "int", "0"
		case "false", "true":
			typeName, result = "bool", v.Name
		case "Errno": // TODO: another heuristic, for go.std.syscall only though
			typeName, result = "uintptr", "0"
		case "Signal": // TODO: another heuristic, for go.std.syscall only though
			typeName, result = "int16", "0" // int16 forces "int()" conversion, which Go requires of "type Signal int"!
		}
		if v.Obj != nil {
			switch spec := v.Obj.Decl.(type) {
			case *ValueSpec:
				if len(spec.Values) == 0 {
					typeName, result = "int", "1" // TODO: probably a good guess for now
				} else {
					typeName, result = evalConstExpr(spec.Values[0])
				}
			case *TypeSpec:
				typeName = evalConstType(spec)
			}
		}
	case *CallExpr:
		typeName, result = evalConstExpr(v.Fun)
	}
	return
}

func determineConstExprType(val Expr) (typeName string) {
	switch v := val.(type) {
	case *BasicLit:
		switch v.Kind {
		case token.STRING:
			typeName = "string"
		case token.INT:
			typeName = fitInt(v.Value)
		case token.FLOAT:
			typeName = "float64"
		case token.CHAR:
			typeName = "rune"
		}
	default:
		typeName, _ = evalConstExpr(val)
	}
	return
}

func determineType(name string, valType, val Expr) (cl, gl string) {
	switch name {
	case "InvalidHandle": // TODO: uintptr on Windows; not found elsewhere
		return "Number", "uint64(%s)"
	}
	typeName := ""
	innerPromotion := "%s"
	if valType == nil {
		typeName = determineConstExprType(val)
	} else {
		ident, ok := valType.(*Ident)
		if !ok {
			return
		}
		valObj := ident.Obj
		if valObj != nil {
			if valObj.Decl != nil {
				ts, ok := valObj.Decl.(*TypeSpec)
				if !ok {
					return
				}
				if ts.Name == nil {
					return
				}
				if id, ok := ts.Type.(*Ident); ok {
					typeName = id.Name
				}
				innerPromotion = typeName + "(%s)"
			}
		} else {
			typeName = ident.Name
		}
	}
	if typeName == "" {
		return
	}
	ti := TypeInfoForGoName(typeName)
	if ti == nil || ti.ArgClojureArgType() == "" || ti.PromoteType() == "" {
		if typeName == "Errno" { // Special-case syscall/zerrors_*.go
			return "Number", "uint64(%s)"
		}
		fmt.Fprintf(os.Stderr, "walk.go/determineType: bad type `%s' at %s\n", typeName, godb.WhereAt(val.Pos()))
		return "", ""
	}
	return ti.ArgClojureArgType(), fmt.Sprintf(ti.PromoteType(), innerPromotion)
}

// Constants are currently emitted while walking the packages. Unlike with variables, where the types are not needed,
// this code seemingly must determine the type of a constant so as to give the Joker wrapper the appropriate type (and
// that is the straightforward way to handle this).
//
// In Go, constants can be explicitly typed, implicitly typed via the constant expressions to which they're assigned, or
// untyped via untyped constant expressions.
//
// Further, Go allows those expressions to refer to constants in other packages, to invoke constructors (say, for simple
// named types like "Type Foo Int") in other packages (as in "const x = Foo(22)", which gives x the type Foo and the
// value 22), and other such things.
//
// Since this code currently makes a complete determination of a constant's type during package walking, it can't count
// on being able to determine the type of anything in another package in order to infer the type that will be given to
// the constant.
//
// Even when all the info is available, this code currently does not attempt to properly evaluate a constant expression
// in order to assure that (for examples) "1 << 30" is "int", "1 << 31" could be "uint" (need to check that), "1 << 32"
// is "int64", "1 << 63" might be "uint64", and so on.
//
// Instead, this code uses some heuristics, including known names of things in Go 1.12's std library, to guess mostly
// correctly, erring on the side of being conservative (which usually means constants that could fit in an "Int" are
// instead a "BigInt").
//
// An attempt was made to change one constant to "variable style" in order to try to eliminate the need for determining
// the type info, via e.g.:
//
//   var EXFULL = syscall.EXFULL
//   var EXFULL_ *GoVar = &GoVar{Value: EXFULL}
//
// That yielded a GoVar[syscall.Errno] type that (int) couldn't convert (probably because it hasn't been special-cased
// to handle GoObject types).
//
// Changing that first line to
//
//   var EXFULL = int(syscall.EXFULL)
//
// solved the problem (EXFULL printed out as an integer, though couldn't be simply, say, added to another integer due to
// being a GoVar[int]), but obviously brings things back to needing to know the type.
//
// Going back to that first approach, and adding this (quick-kludge) code to the procInt function in procs.go:
//
//   	case GoObject:
//		return Int{I: int(obj.O.(syscall.Errno))}
//	case *GoVar:
//		return Int{I: int(obj.Value.(syscall.Errno))}
//
// This allowed "(int EXFULL)" (also "(int (deref EXFULL))", i.e. (int <GoObject[syscall.Errno]>)) to work in that it
// evaluates to a Joker object of type "Int".
//
// Though the kludge (special-casing syscall.Errno) above can be automated away, it doesn't seem like having to always
// wrap such constants in a converter is a helpful requirement.
//
// This isn't just an issue with a named type wrapping a builtin type; even this didn't allow direct referencing of E as
// "Double":
//
//    var E = math.E
//    var E_ *GoVar = &GoVar{Value: E}
//
// It's possible (and perhaps desirable anyway?) that Joker could automatically cast (convert) all GoVar and GoObject
// values to their builtin equivalents, which might allow this all to make more sense.
//
// But it might actually be less work to move the determination of a constant's type to the code-generation phase (so it
// has access to all the packages on which constant expressions might depend) and fully evaluate constant expressions to
// faithfully determine not only their types, but their values as well, and just use those (so, no need to import
// dependent packages).
//
// There might even be an existing Go package to do some of the heavy lifting in that direction. In any case, the result
// would be a lot cleaner and clearer than having "constants" wrapped in GoObject's or GoVar's.
func processConstantSpec(gf *godb.GoFile, pkg string, name *Ident, valType Expr, val Expr, docString string) bool {
	defer func() {
		if x := recover(); x != nil {
			// fmt.Fprintf(os.Stderr, "(Panic due to: %s: %+v)\n", godb.WhereAt(name.Pos()), x)
		}
	}()

	clName := name.Name
	localName := gf.Package.BaseName + "." + name.Name
	fullName := pkg + "." + name.Name

	if typeAndValue, found := typeCheckerInfo.Types[val]; found {
		fmt.Printf("walk.go/processConstantSpec: %s.%s == %s (type %s)\n", pkg, name, typeAndValue.Value, typeAndValue.Type)
	}

	if c, ok := GoConstants[fullName]; ok {
		fmt.Fprintf(os.Stderr, "WARNING: constant %s found at %s and now again at %s\n",
			localName, godb.WhereAt(c.Name.NamePos), godb.WhereAt(name.NamePos))
	}

	switch name.Name {
	case "Int", "String", "Boolean":
		clName += "-renamed" // TODO: is there a better solution possible?
	}

	valTypeString, promoteType := determineType(name.Name, valType, val)
	if WalkDump || (godb.Verbose && valTypeString == "**FOO**") { // or "**FOO**" to quickly disable this
		fmt.Printf("Constant %s at %s:\n", name, godb.WhereAt(name.Pos()))
		if valType != nil {
			fmt.Printf("  valType at %s:\n", godb.WhereAt(valType.Pos()))
			Print(godb.Fset, valType)
		}
		if val != nil {
			fmt.Printf("  val at %s:\n", godb.WhereAt(val.Pos()))
			Print(godb.Fset, val)
		}
	}
	if valTypeString == "" {
		return false
	}

	goCode := fmt.Sprintf(promoteType, localName)

	// Note: :tag value is a string to avoid conflict with like-named member of namespace
	constantDefInfo := map[string]string{
		"DocString":       docString,
		"ValueTypeString": valTypeString,
		"GoCode":          goCode,
		"ClojureName":     clName,
	}

	buf := new(bytes.Buffer)
	Templates.ExecuteTemplate(buf, "clojure-constant-def.tmpl", constantDefInfo)

	gt := &ConstantInfo{name, gf, buf.String()}
	GoConstants[fullName] = gt
	NumGeneratedConstants++

	return true
}

// Note that the 'val' argument isn't used (except when dumping info)
// as it isn't needed to determine the type of a variable, since the
// type isn't needed for code generation for variables -- just for
// constants.
func processVariableSpec(gf *godb.GoFile, pkg string, name *Ident, valType Expr, val Expr, docString string) bool {
	clName := name.Name
	localName := gf.Package.BaseName + "." + name.Name
	fullName := pkg + "." + name.Name

	if c, ok := GoVariables[fullName]; ok {
		fmt.Fprintf(os.Stderr, "WARNING: variable %s found at %s and now again at %s\n",
			localName, godb.WhereAt(c.Name.NamePos), godb.WhereAt(name.NamePos))
	}

	switch name.Name {
	case "Int", "String", "Boolean":
		clName += "-renamed" // TODO: is there a better solution possible?
	}

	if WalkDump {
		fmt.Printf("Variable %s at %s:\n", name, godb.WhereAt(name.Pos()))
		if valType != nil {
			fmt.Printf("  valType at %s:\n", godb.WhereAt(valType.Pos()))
			Print(godb.Fset, valType)
		}
		if val != nil {
			fmt.Printf("  val at %s:\n", godb.WhereAt(val.Pos()))
			Print(godb.Fset, val)
		}
	}

	// Note: :tag value is a string to avoid conflict with like-named member of namespace
	variableDefInfo := map[string]string{
		"DocString":   docString,
		"LocalName":   localName,
		"ClojureName": clName,
	}

	buf := new(bytes.Buffer)
	Templates.ExecuteTemplate(buf, "clojure-variable-def.tmpl", variableDefInfo)

	gt := &VariableInfo{name, gf, buf.String()}
	GoVariables[fullName] = gt
	NumGeneratedVariables++

	return true
}

func what(constant bool) string {
	if constant {
		return "Constant"
	}
	return "Variable"
}

func processValueSpecs(gf *godb.GoFile, pkg string, tss []Spec, parentDoc *CommentGroup, constant bool) {
	var previousVal, previousValType Expr
	for ix, spec := range tss {
		ts := spec.(*ValueSpec)
		for jx, valName := range ts.Names {
			valType := ts.Type
			var val Expr
			if ts.Values != nil {
				if jx >= len(ts.Values) {
					// This seems crazy (more names receiving values than values??) until one
					// investigates the single case that hits this, os/executable_procfs.go, which
					// does something like "var a, b = func() (bool, bool) { ... }()", i.e. the
					// names receive the values returned by the function.
					val = nil
				} else {
					val = ts.Values[jx]
				}
			}

			if val == nil {
				val = previousVal
			}
			if valType == nil {
				valType = previousValType
			}

			if constant {
				previousVal = val
				previousValType = valType
			}

			if !IsExported(valName.Name) {
				continue
			}
			if constant {
				NumConstants++
			} else {
				NumVariables++
			}

			if WalkDump {
				fmt.Printf("%s #%d of spec #%d %s at %s:\n", what(constant), jx, ix, valName, godb.WhereAt(valName.NamePos))
				if valType != nil {
					fmt.Printf("  valType:\n")
					Print(godb.Fset, valType)
				}
				if val != nil {
					fmt.Printf("  val:\n")
					Print(godb.Fset, val)
				}
			}
			doc := ts.Doc // Try block comments for this specific decl
			if doc == nil {
				doc = ts.Comment // Use line comments if no preceding block comments are available
			}
			if doc == nil {
				doc = parentDoc // Use 'var'/'const' statement block comments as last resort
			}
			docString := genutils.CommentGroupInQuotes(doc, "", "", "", "")
			if constant {
				processConstantSpec(gf, pkg, valName, valType, val, docString)
			} else {
				processVariableSpec(gf, pkg, valName, valType, val, docString)
			}
		}
	}
}

func declFunc(gf *godb.GoFile, pkgDirUnix string, f *File, v *FuncDecl) {
	if !IsExported(v.Name.Name) {
		return // Skipping non-exported functions
	}
	if v.Recv != nil {
		for _, r := range v.Recv.List {
			if !astutils.IsExportedType(&r.Type) {
				return // Publishable receivers must operate on public types
			}
		}
		NumReceivers++
	} else {
		NumStandalones++
	}
	NumFunctions++
	processFuncDecl(gf, pkgDirUnix, f, v)
}

func declType(gf *godb.GoFile, pkgDirUnix string, f *File, v *GenDecl) {
	processTypeDecls(gf, pkgDirUnix, v.Specs, v.Doc)
}

func declTypesForTypes(gf *godb.GoFile, pkgDirUnix string, f *File, v *GenDecl) {
	processTypesForTypeDecls(gf, pkgDirUnix, v.Specs, v.Doc)
}

func declValueSpecForTypes(gf *godb.GoFile, pkgDirUnix string, specs []Spec, doc *CommentGroup) {
	processValueSpecsForTypes(gf, pkgDirUnix, specs, doc)
}

func declConstSpec(gf *godb.GoFile, pkgDirUnix string, specs []Spec, doc *CommentGroup) {
	processValueSpecs(gf, pkgDirUnix, specs, doc, true)
}

func declVarSpec(gf *godb.GoFile, pkgDirUnix string, specs []Spec, doc *CommentGroup) {
	processValueSpecs(gf, pkgDirUnix, specs, doc, false)
}

type fileDeclFuncs struct {
	FuncDecl  func(*godb.GoFile, string, *File, *FuncDecl)
	TypeDecl  func(*godb.GoFile, string, *File, *GenDecl)
	ConstDecl func(*godb.GoFile, string, []Spec, *CommentGroup)
	VarDecl   func(*godb.GoFile, string, []Spec, *CommentGroup)
}

func processDecls(gf *godb.GoFile, pkgDirUnix string, f *File, declFuncs fileDeclFuncs) {
	for _, s := range f.Decls {
		switch v := s.(type) {
		case *FuncDecl:
			if declFuncs.FuncDecl != nil {
				declFuncs.FuncDecl(gf, pkgDirUnix, f, v)
			}
		case *GenDecl:
			switch v.Tok {
			case token.TYPE:
				if declFuncs.TypeDecl != nil {
					declFuncs.TypeDecl(gf, pkgDirUnix, f, v)
				}
			case token.CONST:
				if declFuncs.ConstDecl != nil {
					declFuncs.ConstDecl(gf, pkgDirUnix, v.Specs, v.Doc)
				}
			case token.VAR:
				if declFuncs.VarDecl != nil {
					declFuncs.VarDecl(gf, pkgDirUnix, v.Specs, v.Doc)
				}
			case token.IMPORT: // Ignore these
			default:
				panic(fmt.Sprintf("unrecognized token %s at: %s", v.Tok.String(), godb.WhereAt(v.Pos())))
			}
		}
	}
}

type phaseFunc func(*godb.GoFile, string, *File)

func phaseTypeDefs(gf *godb.GoFile, pkgDirUnix string, f *File) {
	processDecls(gf, pkgDirUnix, f, fileDeclFuncs{
		FuncDecl:  nil,
		TypeDecl:  declType,
		ConstDecl: nil,
		VarDecl:   nil,
	})
}

func phaseTypeRefs(gf *godb.GoFile, pkgDirUnix string, f *File) {
	processDecls(gf, pkgDirUnix, f, fileDeclFuncs{
		FuncDecl:  declFuncForTypes,
		TypeDecl:  declTypesForTypes,
		ConstDecl: declValueSpecForTypes,
		VarDecl:   declValueSpecForTypes,
	})
}

func phaseOtherDecls(gf *godb.GoFile, pkgDirUnix string, f *File) {
	processDecls(gf, pkgDirUnix, f, fileDeclFuncs{
		FuncDecl:  declFunc,
		TypeDecl:  nil,
		ConstDecl: declConstSpec,
		VarDecl:   declVarSpec,
	})
}

func processPackage(rootUnix, pkgDirUnix, nsRoot string, p *Package, fn phaseFunc) {
	for path, f := range p.Files {
		goFilePathUnix := TrimPrefix(filepath.ToSlash(path), rootUnix+"/")
		gf := godb.GoFilesRelative[goFilePathUnix]
		if gf == nil {
			panic(fmt.Sprintf("cannot find GoFile object for %s", goFilePathUnix))
		}

		fn(gf, pkgDirUnix, f)
	}
}

func processDir(rootNative, pathNative paths.NativePath, nsRoot, importMeRoot string) error {
	pkgDirNative, ok := pathNative.RelativeTo(rootNative)
	if !ok {
		panic(fmt.Sprintf("%s is not relative to %s", pathNative, rootNative))
	}
	pkgDirUnix := pkgDirNative.ToUnix()
	if godb.Verbose {
		genutils.AddSortedStdout(fmt.Sprintf("Processing %s:\n", pkgDirUnix))
	}
	importMe := path.Join(importMeRoot, pkgDirUnix.String())

	pkgs, err := parser.ParseDir(godb.Fset, pathNative.String(),
		// Walk only *.go files that meet default (target) build constraints, e.g. per "// build ..."
		func(info os.FileInfo) bool {
			if HasSuffix(info.Name(), "_test.go") {
				if godb.Verbose {
					genutils.AddSortedStdout(fmt.Sprintf("Ignoring test code in %s\n", info.Name()))
				}
				return false
			}
			b, e := build.Default.MatchFile(pathNative.String(), info.Name())
			if godb.Verbose {
				genutils.AddSortedStdout(fmt.Sprintf("Matchfile(%s) => %v %v\n",
					pathNative.Join(info.Name()).ToUnix(),
					b, e))
			}
			return b && e == nil
		},
		parser.ParseComments|parser.AllErrors)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	found := false
	for pkgBaseName, pkg := range pkgs {
		// fmt.Println(pkg.ID, pkg.GoFiles)
		// info := &types.Info{}
		// checkedPkg, err := typesConf.Check("hey/what", godb.Fset, pkg.GoFiles, info)

		if pkgBaseName != pathNative.Base() {
			if godb.Verbose {
				genutils.AddSortedStdout(fmt.Sprintf("NOTICE: Package %s is defined in %s -- ignored due to name mismatch\n",
					pkgBaseName, pathNative))
			}
		} else {
			if found {
				panic("whaaa??")
			}
			// Cannot currently do this, as public constants generated via "_ Something = iota" are omitted:
			// FilterPackage(v, IsExported)
			godb.RegisterPackage(rootNative.ToUnix(), pkgDirUnix, nsRoot, importMe, pkg)
			found = true
		}
	}

	return nil
}

var excludeDirs = map[string]bool{
	"builtin":  true,
	"cmd":      true,
	"internal": true, // look into this later?
	"testdata": true,
	"vendor":   true,
}

func LegitimateImport(p string) bool {
	if p == "C" {
		return false
	}
	elements := Split(p, "/")
	for _, e := range elements {
		if excludeDirs[e] {
			return false
		}
	}
	return true
}

func walkDir(fsRoot paths.NativePath, nsRoot, importMeRoot string) error {
	target, err := fsRoot.EvalSymlinks()
	Check(err)

	err = target.Walk(
		func(path paths.NativePath, info os.FileInfo, err error) error {
			rel := ReplaceAll(path.String(), target.String(), fsRoot.String())
			relNative := paths.NewNativePath(rel)
			relUnix := relNative.ToUnix()
			if err != nil {
				genutils.EndSortedStdout()
				fmt.Fprintf(os.Stderr, "Skipping %s due to: %v\n", relUnix, err)
				return err
			}
			if relNative == fsRoot {
				return nil // skip (implicit) "."
			}
			if excludeDirs[relUnix.Base()] {
				if godb.Verbose {
					genutils.AddSortedStdout(fmt.Sprintf("Excluding %s\n", relUnix))
				}
				return paths.SkipDir
			}
			if info.IsDir() {
				return processDir(fsRoot, relNative, nsRoot, importMeRoot)
			}
			return nil // not a directory
		})

	if err != nil {
		genutils.EndSortedStdout()
		fmt.Fprintf(os.Stderr, "Error while walking %s: %v\n", fsRoot, err)
		return err
	}

	return err
}

type dirToWalk struct {
	srcDir   paths.NativePath
	fsRoot   paths.NativePath
	nsRoot   string
	importMe string
}

var dirsToWalk []dirToWalk

func AddWalkDir(srcDir, fsRoot paths.NativePath, nsRoot, importMe string) {
	dirsToWalk = append(dirsToWalk, dirToWalk{srcDir, fsRoot, nsRoot, importMe})
}

var typeCheckerConfig *types.Config
var typeCheckerInfo *types.Info

type importerFunc func(path string) (*types.Package, error)

func (f importerFunc) Import(path string) (*types.Package, error) {
	return f(path)
}

func myImporter(path string) (*types.Package, error) {
	if path == "unsafe" {
		return types.Unsafe, nil
	}
	pkg := godb.GetPackageInfo(path)
	if pkg == nil {
		return nil, nil // TODO: Something better when package not found?
	}
	files := []*File{}
	for _, f := range pkg.Files {
		files = append(files, f)
	}
	return typeCheckerConfig.Check(path, godb.Fset, files, typeCheckerInfo)
}

func WalkAllDirs() (error, paths.NativePath) {
	var phases = []phaseFunc{
		phaseTypeDefs,
		phaseTypeRefs,
		phaseOtherDecls,
	}

	genutils.StartSortedStdout()
	defer func() {
		genutils.EndSortedStdout()
	}()

	for _, d := range dirsToWalk {
		err := walkDir(d.fsRoot, d.nsRoot, d.importMe)
		if err != nil {
			return err, d.srcDir
		}
	}

	//	importer := importerFunc(myImporter)

	typeCheckerConfig = &types.Config{
		IgnoreFuncBodies: true,
		FakeImportC:      true,
		Importer:         importer.Default(),
	}
	typeCheckerInfo = &types.Info{
		Types: map[Expr]types.TypeAndValue{},
		Defs:  map[*Ident]types.Object{},
	}

	if _, err := myImporter("net"); err != nil {
		fmt.Fprintf(os.Stderr, "walk.go/WalkAllDirs(): Failed to check %q: %s\n", "net", err)
	}

	for _, wp := range godb.PackagesAsDiscovered {
		initPackage(wp.Root.String(), wp.Dir.String(), wp.NsRoot, wp.Pkg)
	}

	for _, phase := range phases {
		for _, wp := range godb.PackagesAsDiscovered {
			processPackage(wp.Root.String(), wp.Dir.String(), wp.NsRoot, wp.Pkg, phase)
		}
	}

	return nil, paths.NewNativePath("")
}

func findApis(src paths.NativePath) (apis map[string]struct{}) {
	start := getCPU()
	defer func() {
		end := getCPU()
		if godb.Verbose && !noTimeAndVersion {
			fmt.Printf("findApis() took %d ns.\n", end-start)
		}
	}()

	apis = map[string]struct{}{}

	var fset = token.NewFileSet()

	target, err := src.ToNative().EvalSymlinks()
	Check(err)

	pkgs, err := parser.ParseDir(fset, target.String(), nil, 0)
	Check(err)

	var pkg *Package
	for k, v := range pkgs {
		if k != "core" {
			panic(fmt.Sprintf("Expected only package 'core', found '%s'", k))
		}
		pkg = v
	}

	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			switch o := d.(type) {
			case *FuncDecl:
				if o.Recv == nil {
					if IsExported(o.Name.Name) {
						apis[o.Name.Name] = struct{}{}
					}
				}
			}
		}
	}

	return
}

// Determine the runtime API name for a function call, given a choice
// of prefixes (core and namespace-based) and the type name.  Ensure
// the resulting API has been code-generated or already exists in
// package core, wrapping it in an ABEND if not, and return the
// resulting wrap or, if no errors, the original string (which
// generate-std.joke will use to reconstitute the same API name).
func assertRuntime(prefix, nsPrefix, s string) string {
	runtime := s
	if ix := Index(s, "("); ix >= 0 {
		runtime = runtime[0:ix]
	}
	if ix := Index(runtime, "/"); ix >= 0 {
		ns := runtime[0 : ix+1]
		runtime = ns + nsPrefix + runtime[ix+1:]
	} else {
		runtime = prefix + runtime
	}
	if Contains(runtime, "ABEND") {
		return s
	}
	if _, found := definedApis[runtime]; !found {
		return fmt.Sprintf("ABEND707(API '%s' is unimplemented: %s)", runtime, s)
	}
	return s
}

// Determines (and validates) the API to call (in context), given the
// full Clojure typename (e.g. "go.std.something/Foo" or
// "arrayOfByte"), the import base name, and choice of prefixes.
func determineRuntime(prefix, nsPrefix, imp, clType string) string {
	var runtime, call string
	if ix := Index(clType, "/"); ix >= 0 {
		runtime = clType[0:ix+1] + nsPrefix + clType[ix+1:]
		if imp != "" {
			imp += "."
		}
		call = imp + nsPrefix + clType[ix+1:]
	} else {
		runtime = prefix + clType
		call = runtime
	}
	if Contains(runtime, "ABEND") {
		return runtime
	}
	if _, found := definedApis[runtime]; !found {
		return fmt.Sprintf("ABEND707(API '%s' is unimplemented: %s)", runtime, clType)
	}
	return call
}
