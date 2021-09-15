package main

import (
	"fmt"
	"github.com/candid82/joker/tools/gostd/astutils"
	"github.com/candid82/joker/tools/gostd/genutils"
	"github.com/candid82/joker/tools/gostd/godb"
	"github.com/candid82/joker/tools/gostd/gtypes"
	"github.com/candid82/joker/tools/gostd/imports"
	"github.com/candid82/joker/tools/gostd/jtypes"
	. "go/ast"
	"go/token"
	"go/types"
	"os"
	"sort"
)

type TypeInfo interface {
	ArgClojureType() string       // Can convert this type to a Go function arg with my type
	ArgFromClojureObject() string /// Append this to Clojure object to extract value of my type
	ArgExtractFunc() string       // Call Extract<this>() for arg with my type
	ArgClojureArgType() string    // Clojure argument type for a Go function arg with my type
	ConvertFromClojure() string
	ConvertFromMap() string  // Pattern to convert a map %s key %s to this type
	AsClojureObject() string // Pattern to convert this type to a normal Clojure type, or empty string to simply wrap in a GoObject
	ClojureName() string
	ClojureEffectiveName() string
	ClojureNameDocForType(*types.Package) string
	ClojurePattern() string
	ClojureBaseName() string
	ClojureTypeInfo() *jtypes.Info
	RequiredImports() *imports.Imports
	GoName() string
	GoNameDocForType(*types.Package) string // Relative to pkg
	GoPackage() string
	GoPattern() string
	GoBaseName() string
	GoEffectiveBaseName() string // Substitutes what actually works in generated Go code (interface{} instead of Arbitrary if in unsafe pkg)
	GoTypeExpr() Expr            // The actual (Go) type (if any)
	GoTypeInfo() *gtypes.Info
	TypeSpec() *TypeSpec // Definition, if any, of named type
	UnderlyingTypeInfo() TypeInfo
	UnderlyingType() Expr // nil if not a declared type
	UnderlyingTypeSpec() *TypeSpec
	GoFile() *godb.GoFile
	DefPos() token.Pos
	Specificity() uint // ConcreteType, else # of methods defined for interface{} (abstract) type
	PromoteType() string
	TypeMappingsName() string
	Doc() string
	NilPattern() string
	IsCustom() bool      // Whether this is defined by the codebase vs either builtin or so derived
	IsUnsupported() bool // Is this unsupported?
	IsNullable() bool    // Can an instance of the type == nil (e.g. 'error' type)?
	IsExported() bool
	IsSwitchable() bool      // Can (Go) name be used in a 'case' statement or type assertion?
	IsAddressable() bool     // Is "&instance" going to pass muster, even with 'go vet'?
	IsPassedByAddress() bool // Excludes builtins, some complex, and interface{} types
	IsArbitraryType() bool   // Is unsafe.ArbitraryType, which gets treated as interface{}
	IsCtorable() bool        // Whether a ctor for this type can (and will) be created
}

type TypesMap map[string]TypeInfo

// Maps type-defining Expr or full names to exactly one struct describing that type.
var typesByExpr = map[Expr]TypeInfo{}
var typesByGoName = TypesMap{}
var typesByClojureName = TypesMap{}
var typesByGoTypeName = TypesMap{}

var typesByGoType = map[types.Type]TypeInfo{}

const ConcreteType = gtypes.Concrete

type typeInfo struct {
	jti             *jtypes.Info
	gti             *gtypes.Info
	requiredImports *imports.Imports
	who             string // who made me
}

func RegisterTypeDecl(ts *TypeSpec, gf *godb.GoFile, pkg string, parentDoc *CommentGroup) {
	name := ts.Name.Name
	goTypeName := pkg + "." + name

	if WalkDump {
		fmt.Printf("Type %s at %s:\n", goTypeName, godb.WhereAt(ts.Pos()))
		Print(godb.Fset, ts)
	}

	gtiVec := gtypes.Define(ts, gf, parentDoc)

	for _, gti := range gtiVec {

		jti := jtypes.Define(ts, gti.Expr)

		ti := &typeInfo{
			jti:             jti,
			gti:             gti,
			requiredImports: &imports.Imports{},
			who:             "RegisterTypeDecl",
		}

		typesByGoName[ti.GoName()] = ti
		typesByClojureName[ti.ClojureName()] = ti
		typesByGoTypeName[ti.GoTypeName()] = ti
		typesByGoType[ti.GoType()] = ti

		if IsExported(name) {
			NumTypes++
			if ti.IsCtorable() {
				NumCtableTypes++
			}
		}

		ClojureCode[pkg].InitTypes[ti] = struct{}{}
		GoCode[pkg].InitTypes[ti] = struct{}{}
	}
}

func RegisterAllSubtypes(e Expr) {
	if e == nil {
		return
	}

	switch v := e.(type) {
	case *Ident:
		return
	case *StarExpr:
		RegisterAllSubtypes(v.X)
	case *ArrayType:
		//		RegisterAllSubtypes(v.Len)
		RegisterAllSubtypes(v.Elt)
	case *InterfaceType:
		for _, f := range astutils.FlattenFieldList(v.Methods) {
			if f.Name == nil || IsExported(f.Name.Name) {
				RegisterAllSubtypes(f.Field.Type)
			}
		}
	case *MapType:
		// RegisterAllSubtypes(v.Key)
		// RegisterAllSubtypes(v.Value)
		return
	case *SelectorExpr:
	case *ChanType:
		RegisterAllSubtypes(v.Value)
	case *StructType:
		for _, f := range astutils.FlattenFieldList(v.Fields) {
			if f.Name == nil || IsExported(f.Name.Name) {
				RegisterAllSubtypes(f.Field.Type)
			}
		}
	case *FuncType:
		return
	case *Ellipsis:
		return
	default:
		return
	}

	TypeInfoForExpr(e)
}

func TypeInfoForExpr(e Expr) TypeInfo {
	if ti, found := typesByExpr[e]; found {
		return ti
	}

	gti := gtypes.InfoForExpr(e)
	if gti == nil {
		return nil
	}
	jti := jtypes.InfoForExpr(e, gti.GoType)

	if ti, found := typesByGoName[gti.FullName]; found {
		if _, ok := typesByClojureName[jti.FullName]; !ok {
			panic(fmt.Sprintf("types.go/TypeInfoForExpr: have typesByGoName[%s] but not typesByClojureName[%s] for %s at %s\n", gti.FullName, jti.FullName, astutils.ExprToString(e), godb.WhereAt(e.Pos())))
			//			typesByClojureName[jti.FullName] = ti
		}
		if _, ok := typesByGoType[gti.GoType]; !ok {
			fmt.Fprintf(os.Stderr, "types.go/TypeInfoForExpr: have typesByGoName[%s] but not typesByGoType[%v] for %s at %s (alias type?)\n", gti.FullName, gti.GoType, astutils.ExprToString(e), godb.WhereAt(e.Pos()))
		}
		return ti
	}
	if _, ok := typesByClojureName[jti.FullName]; ok && jti.FullName != "GoObject" {
		// if inf := jtypes.InfoForName(jti.FullName); inf == nil {
		// 	panic(fmt.Sprintf("types.go/TypeInfoForExpr: have typesByClojureName[%s] but not typesByGoName[%s] for %s at %s\n", jti.FullName, gti.FullName, astutils.ExprToString(e), godb.WhereAt(e.Pos())))
		// }
		if _, ok := typesByGoType[gti.GoType]; ok {
			panic(fmt.Sprintf("types.go/TypeInfoForExpr: have typesByGoType[%v] but not typesByGoName[%s] for %s at %s\n", gti.GoType, gti.FullName, astutils.ExprToString(e), godb.WhereAt(e.Pos())))
		}
	}

	ti := &typeInfo{
		gti:             gti,
		jti:             jti,
		requiredImports: &imports.Imports{},
		who:             "TypeInfoForExpr",
	}

	//	fmt.Printf("types.go/TypeInfoForExpr: @%p; gti: %s == @%p %+v; jti: %s == @%p %+v; at %s\n", ti, ti.GoName(), gti, gti, ti.ClojureName(), ti, ti, godb.WhereAt(e.Pos()))

	typesByExpr[e] = ti
	typesByGoName[gti.FullName] = ti
	typesByClojureName[jti.FullName] = ti
	typesByGoTypeName[gti.TypeName] = ti
	typesByGoType[gti.GoType] = ti

	return ti
}

func TypeInfoForGoName(goName string) TypeInfo {
	if ti, found := typesByGoName[goName]; found {
		return ti
	}

	gti := gtypes.InfoForName(goName)
	if gti == nil {
		panic(fmt.Sprintf("cannot find `%s' in gtypes", goName))
	}

	jti := jtypes.InfoForGoTypeName(goName)
	if jti == nil {
		panic(fmt.Sprintf("cannot find `%s' in jtypes", goName))
	}

	ti := &typeInfo{
		gti:             gti,
		jti:             jti,
		requiredImports: &imports.Imports{},
		who:             "TypeInfoForGoName",
	}

	typesByGoName[gti.FullName] = ti
	typesByClojureName[jti.FullName] = ti
	typesByGoTypeName[gti.TypeName] = ti
	typesByGoType[gti.GoType] = ti

	return ti
}

func TypeInfoForType(ty types.Type) TypeInfo {
	if ti, found := typesByGoType[ty]; found {
		return ti
	}

	gti := gtypes.InfoForType(ty)
	if gti == nil {
		panic(fmt.Sprintf("cannot find `%s' in gtypes", ty.String()))
	}
	typeName := gti.TypeName

	jti := jtypes.InfoForGoType(ty)
	if jti == nil {
		jti = jtypes.InfoForGoTypeName(typeName)
	}
	if jti == nil {
		panic(fmt.Sprintf("cannot find `%s' in jtypes", typeName))
	}

	ti := &typeInfo{
		gti:             gti,
		jti:             jti,
		requiredImports: &imports.Imports{},
		who:             "TypeInfoForType",
	}

	typesByGoName[gti.FullName] = ti
	typesByClojureName[jti.FullName] = ti
	typesByGoTypeName[typeName] = ti
	typesByGoType[gti.GoType] = ti

	return ti
}

func StringForExpr(e Expr) string {
	if e == nil {
		return "-"
	}
	t := TypeInfoForExpr(e)
	if t != nil {
		return t.GoName()
	}
	return fmt.Sprintf("%T", e)
}

func conversions(e Expr) (fromClojure, fromMap string) {
	switch v := e.(type) {
	case *Ident:
		//		fmt.Fprintf(os.Stderr, "conversions(Ident:%+v)\n", v)
		if ti := TypeInfoForGoName(v.Name); ti != nil {
			if ti.IsCustom() {
				uti := TypeInfoForExpr(ti.TypeSpec().Type)
				if uti.ConvertFromClojure() != "" {
					fromClojure = fmt.Sprintf("%s.%s(%s)", ti.GoPackage(), ti.GoBaseName(), uti.ConvertFromClojure())
				}
				if uti.ConvertFromMap() != "" {
					fromMap = fmt.Sprintf("%s.%s(%s)", ti.GoPackage(), ti.GoBaseName(), uti.ConvertFromMap())
				}
			}
		}
	case *InterfaceType:
		if !v.Incomplete && len(v.Methods.List) == 0 {
			fromMap = "FieldAsGoObject(%s, %s)"
			fromClojure = "ObjectAsGoObject(%s, %s)"
		}
	default:
	}
	return
}

func SortedTypeInfoMap(m TypesMap, f func(string, TypeInfo)) {
	var keys []string
	for k, _ := range m {
		if k[0] == '*' {
			keys = append(keys, k[1:]+"*")
		} else {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		if k[len(k)-1] == '*' {
			k = "*" + k[0:len(k)-1]
		}
		f(k, m[k])
	}
}

func (ti typeInfo) ArgClojureType() string {
	return ti.jti.ArgClojureType
}

func (ti typeInfo) ArgFromClojureObject() string {
	return ti.jti.ArgFromClojureObject
}

func (ti typeInfo) ArgExtractFunc() string {
	return ti.jti.ArgExtractFunc
}

func (ti typeInfo) ArgClojureArgType() string {
	return ti.jti.ArgClojureArgType
}

func (ti typeInfo) ConvertFromClojure() string {
	return ti.jti.ConvertFromClojure
}

func (ti typeInfo) ConvertFromMap() string {
	return ti.jti.ConvertFromMap
}

func (ti typeInfo) AsClojureObject() string {
	return ti.jti.AsClojureObject
}

func (ti typeInfo) ClojureName() string {
	return ti.jti.FullName
}

func (ti typeInfo) ClojureEffectiveName() string {
	if ti.gti.IsArbitraryType {
		return "GoObject"
	}
	return ti.jti.FullName
}

func (ti typeInfo) ClojureNameDoc(e Expr) string {
	if ti.gti.IsArbitraryType {
		return "GoObject"
	}
	return ti.jti.NameDoc(e)
}

func (ti typeInfo) ClojureNameDocForType(pkg *types.Package) string {
	if ti.gti.IsArbitraryType {
		return "GoObject"
	}
	return ti.jti.NameDocForType(pkg)
}

func (ti typeInfo) ClojurePattern() string {
	return ti.jti.Pattern
}

func (ti typeInfo) ClojureBaseName() string {
	return ti.jti.BaseName
}

func (ti typeInfo) NilPattern() string {
	return ti.gti.NilPattern
}

func (ti typeInfo) ClojureTypeInfo() *jtypes.Info {
	return ti.jti
}

func (ti typeInfo) ClojureWho() string {
	return ti.jti.Who
}

func (ti typeInfo) PromoteType() string {
	return ti.jti.PromoteType
}

func (ti typeInfo) RequiredImports() *imports.Imports {
	return ti.requiredImports
}

func (ti typeInfo) GoName() string {
	return ti.gti.FullName
}

func (ti typeInfo) GoNameDoc(e Expr) string {
	return ti.gti.NameDoc(e)
}

func (ti typeInfo) GoNameDocForType(pkg *types.Package) string {
	return ti.gti.NameDocForType(pkg)
}

func (ti typeInfo) GoPackage() string {
	return ti.gti.Package
}

func (ti typeInfo) GoPattern() string {
	return ti.gti.Pattern
}

func (ti typeInfo) GoBaseName() string {
	return ti.gti.LocalName
}

func (ti typeInfo) GoEffectiveBaseName() string {
	if ti.gti.IsArbitraryType {
		return "interface{}"
	}
	return ti.gti.LocalName
}

func (ti typeInfo) GoType() types.Type {
	return ti.gti.GoType
}

func (ti typeInfo) GoTypeExpr() Expr {
	return ti.gti.Type
}

func (ti typeInfo) GoTypeInfo() *gtypes.Info {
	return ti.gti
}

func (ti typeInfo) GoTypeName() string {
	return ti.gti.TypeName
}

func (ti typeInfo) TypeSpec() *TypeSpec {
	return ti.gti.TypeSpec
}

func (ti typeInfo) UnderlyingTypeInfo() TypeInfo {
	ut := ti.gti.UnderlyingType
	if ut == nil {
		return nil
	}
	if ut.Expr == nil {
		return typesByGoName[ut.FullName]
	}
	return TypeInfoForExpr(ut.Expr)
}

func (ti typeInfo) UnderlyingType() Expr {
	if ut := ti.gti.UnderlyingType; ut != nil {
		return ut.Expr
	}
	return nil
}

func (ti typeInfo) UnderlyingTypeSpec() *TypeSpec {
	if ts := ti.gti.TypeSpec; ts != nil {
		return ts
	}
	if uti := ti.UnderlyingTypeInfo(); uti != nil {
		return uti.UnderlyingTypeSpec()
	}
	return nil
}

func (ti typeInfo) GoFile() *godb.GoFile {
	return ti.gti.File
}

func (ti typeInfo) DefPos() token.Pos {
	return ti.gti.DefPos
}

func (ti typeInfo) Specificity() uint {
	return ti.gti.Specificity
}

func (ti typeInfo) Doc() string {
	return ti.gti.Doc
}

func (ti typeInfo) TypeMappingsName() string {
	if ugt := ti.gti.UnderlyingType; ugt != nil {
		switch ti.gti.Expr.(type) {
		case *ArrayType:
			return "info_ArrayOf_" + fmt.Sprintf(ugt.Pattern, ugt.LocalName)
		case *StarExpr:
			return "info_PtrTo_" + fmt.Sprintf(ugt.Pattern, ugt.LocalName)
		default:
			panic(fmt.Sprintf("unexpected expr %T with underlying type", ti.gti.Expr))
		}
	}
	return "info_" + fmt.Sprintf(ti.GoPattern(), ti.GoBaseName())
}

func (ti typeInfo) Namespace() string {
	return ti.jti.Namespace
}

func (ti typeInfo) IsCustom() bool {
	return ti.TypeSpec() != nil || ti.UnderlyingTypeInfo() != nil
}

func (ti typeInfo) IsUnsupported() bool {
	return ti.gti.IsUnsupported || ti.jti.IsUnsupported
}

func (ti typeInfo) IsNullable() bool {
	return ti.gti.IsNullable
}

func (ti typeInfo) IsExported() bool {
	return ti.gti.IsExported
}

func (ti typeInfo) IsBuiltin() bool {
	return ti.gti.IsBuiltin
}

func (ti typeInfo) IsSwitchable() bool {
	return ti.gti.IsSwitchable
}

func (ti typeInfo) IsAddressable() bool {
	return ti.gti.IsAddressable
}

func (ti typeInfo) IsPassedByAddress() bool {
	return ti.gti.IsPassedByAddress
}

func (ti typeInfo) IsArbitraryType() bool {
	return ti.gti.IsArbitraryType
}

func (ti typeInfo) IsCtorable() bool {
	return ti.gti.IsCtorable
}

func (ti typeInfo) IsReferenced() bool {
	return ti.gti.IsExported || ti.GoName() == "net.conn" || ti.GoName() == "net.*conn"
}

var myAllTypesSorted = []TypeInfo{}

// This establishes the order in which types are matched by 'case' statements in the "big switch" in goswitch.go. Once established,
// new types cannot be discovered/added.
func SortAllTypes() []TypeInfo {
	if len(myAllTypesSorted) > 0 {
		panic("Attempt to sort all types type after having already sorted all types!!")
	}
	for _, ti := range typesByClojureName {
		if ti.GoName() == "[][]*crypto/x509.Certificate XXX DISABLED XXX" {
			fmt.Printf("types.go/SortAllTypes: %s == %+v %+v\n", ti.ClojureName(), ti.GoTypeInfo(), ti.ClojureTypeInfo())
		}
		t := ti.GoTypeInfo()
		if t.IsExported && !t.IsArbitraryType && !t.IsBuiltin {
			myAllTypesSorted = append(myAllTypesSorted, ti.(*typeInfo))
		}
	}
	sort.SliceStable(myAllTypesSorted, func(i, j int) bool {
		i_gti := myAllTypesSorted[i].GoTypeInfo()
		j_gti := myAllTypesSorted[j].GoTypeInfo()
		if iSpecificity, jSpecificity := i_gti.Specificity, j_gti.Specificity; iSpecificity != jSpecificity {
			return iSpecificity > jSpecificity
		}
		return i_gti.FullName < j_gti.FullName
	})
	return myAllTypesSorted
}

func typeKeyForSort(ti TypeInfo) string {
	// Put the pattern (e.g. "%s", "*%s") at the end, to put related types together in a reasonable order.
	return genutils.CombineGoName(ti.GoPackage(), ti.GoBaseName()+ti.GoPattern())
}

func SortedTypeDefinitions(m map[TypeInfo]struct{}, f func(TypeInfo)) {
	var keys []string
	vals := TypesMap{}
	for k, _ := range m {
		if k != nil {
			key := k.GoTypeInfo().FullName
			keys = append(keys, key)
			vals[key] = k
		}
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return typeKeyForSort(vals[keys[i]]) < typeKeyForSort(vals[keys[j]])
	})
	for _, k := range keys {
		f(vals[k])
	}
}

func TypesByGoName() TypesMap {
	return typesByGoName
}

func init() {
	jtypes.ConversionsFn = conversions
}
