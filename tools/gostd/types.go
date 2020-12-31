package main

import (
	"fmt"
	"github.com/candid82/joker/tools/gostd/genutils"
	"github.com/candid82/joker/tools/gostd/godb"
	"github.com/candid82/joker/tools/gostd/gtypes"
	"github.com/candid82/joker/tools/gostd/imports"
	"github.com/candid82/joker/tools/gostd/jtypes"
	. "go/ast"
	"go/token"
	"sort"
)

type TypeInfo interface {
	ArgClojureType() string       // Can convert this type to a Go function arg with my type
	ArgFromClojureObject() string /// Append this to Clojure object to extract value of my type
	ArgExtractFunc() string       // Call Extract<this>() for arg with my type
	ArgClojureArgType() string    // Clojure argument type for a Go function arg with my type
	ConvertFromClojure() string
	ConvertFromMap() string  // Pattern to convert a map %s key %s to this type
	Custom() bool            // Whether this is defined by the codebase vs either builtin or derived
	AsClojureObject() string // Pattern to convert this type to a normal Clojure type, or empty string to simply wrap in a GoObject
	ClojureName() string
	ClojureEffectiveName() string
	ClojureNameDoc(e Expr) string
	ClojurePattern() string
	ClojureBaseName() string
	ClojureTypeInfo() *jtypes.Info
	ClojureWho() string
	RequiredImports() *imports.Imports
	GoName() string
	GoEffectiveName() string // Substitutes what actually works in generated Go code (inteface{} instead of unsafe.Arbitrary)
	GoNameDoc(e Expr) string
	GoPackage() string
	GoPattern() string
	GoBaseName() string
	GoEffectiveBaseName() string // Substitutes what actually works in generated Go code (inteface{} instead of Arbitrary if in unsafe pkg)
	GoTypeInfo() *gtypes.Info
	TypeSpec() *TypeSpec // Definition, if any, of named type
	UnderlyingTypeInfo() TypeInfo
	UnderlyingType() Expr // nil if not a declared type
	GoFile() *godb.GoFile
	DefPos() token.Pos
	Specificity() uint // ConcreteType, else # of methods defined for interface{} (abstract) type
	PromoteType() string
	TypeMappingsName() string
	Doc() string
	IsUnsupported() bool // Is this unsupported?
	IsNullable() bool    // Can an instance of the type == nil (e.g. 'error' type)?
	IsExported() bool
	IsBuiltin() bool
	IsSwitchable() bool      // Can (Go) name be used in a 'case' statement or type assertion?
	IsAddressable() bool     // Is "&instance" going to pass muster, even with 'go vet'?
	IsPassedByAddress() bool // Excludes builtins, some complex, and interface{} types
	IsArbitraryType() bool   // Is unsafe.ArbitraryType, which gets treated as interface{}
}

type TypesMap map[string]TypeInfo

// Maps type-defining Expr or full names to exactly one struct describing that type.
var typesByExpr = map[Expr]TypeInfo{}
var typesByGoName = TypesMap{}
var typesByClojureName = TypesMap{}

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

		if IsExported(name) {
			NumTypes++
			if ti.Specificity() == ConcreteType {
				NumCtableTypes++
			}
		}

		ClojureCode[pkg].InitTypes[ti] = struct{}{}
		GoCode[pkg].InitTypes[ti] = struct{}{}
	}
}

func TypeInfoForExpr(e Expr) TypeInfo {
	if ti, found := typesByExpr[e]; found {
		return ti
	}

	gti := gtypes.InfoForExpr(e)
	jti := jtypes.InfoForExpr(e)

	if ti, found := typesByGoName[gti.FullName]; found {
		if _, ok := typesByClojureName[jti.FullName]; !ok {
			fmt.Printf("types.go/TypeInfoForExpr: have typesByGoName[%s] but not typesByClojureName[%s]\n", gti.FullName, jti.FullName)
			typesByClojureName[jti.FullName] = ti
		}
		return ti
	}
	if _, ok := typesByClojureName[jti.FullName]; ok && jti.FullName != "GoObject" {
		if inf := jtypes.InfoForGoName(jti.FullName); inf == nil {
			fmt.Printf("types.go/TypeInfoForExpr: have typesByClojureName[%s] but not typesByGoName[%s]\n", jti.FullName, gti.FullName)
		}
	}

	ti := &typeInfo{
		gti:             gti,
		jti:             jti,
		requiredImports: &imports.Imports{},
		who:             "TypeInfoForExpr",
	}

	//	fmt.Printf("types.go/TypeInfoForExpr: %s == @%p %+v at %s\n", ti.ClojureName(), ti, ti, godb.WhereAt(e.Pos()))

	typesByExpr[e] = ti
	typesByGoName[gti.FullName] = ti
	typesByClojureName[jti.FullName] = ti

	return ti
}

func TypeInfoForGoName(goName string) TypeInfo {
	if ti, found := typesByGoName[goName]; found {
		return ti
	}

	gti := gtypes.InfoForName(goName)
	if gti == nil {
		return nil // panic(fmt.Sprintf("cannot find `%s' in gtypes", goName))
	}

	jti := jtypes.InfoForGoName(goName)
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
			if ti.Custom() {
				uti := TypeInfoForExpr(ti.TypeSpec().Type)
				if uti.ConvertFromClojure() != "" {
					fromClojure = fmt.Sprintf("%s.%s(%s)", ti.GoPackage(), ti.GoBaseName(), uti.ConvertFromClojure())
				}
				if uti.ConvertFromMap() != "" {
					fromMap = fmt.Sprintf("%s.%s(%s)", ti.GoPackage(), ti.GoBaseName(), uti.ConvertFromMap())
				}
			}
		}
	case *ArrayType:
		//		fmt.Fprintf(os.Stderr, "conversions(ArrayType:%+v)\n", v)
	case *StarExpr:
		//		fmt.Fprintf(os.Stderr, "conversions(StarExpr:%+v)\n", v)
		// ti := TypeInfoForExpr(v.X)
		// if ti.ConvertFromClojure() != "" && ti.ArgClojureArgType() == ti.ArgExtractFunc() {
		// 	fromClojure = "&" + ti.ConvertFromClojure()
		// }
		// if ti.ConvertFromMap() != "" && ti.ArgClojureArgType() == ti.ArgExtractFunc() {
		// 	//			fromMap = "&" + ti.ConvertFromMap()
		// }
	default:
		//		fmt.Fprintf(os.Stderr, "conversions(default:%+v)\n", v)
	}
	return
}

func SortedTypeInfoMap(m TypesMap, f func(k string, v TypeInfo)) {
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

func (ti typeInfo) Custom() bool {
	return ti.TypeSpec() != nil || ti.UnderlyingTypeInfo() != nil
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

func (ti typeInfo) ClojurePattern() string {
	return ti.jti.Pattern
}

func (ti typeInfo) ClojureBaseName() string {
	return ti.jti.BaseName
}

func (ti typeInfo) IsUnsupported() bool {
	return ti.jti.IsUnsupported
}

func (ti typeInfo) ClojureTypeInfo() *jtypes.Info {
	return ti.jti
}

func (ti typeInfo) ClojureWho() string {
	return ti.jti.Who
}

func (ti typeInfo) RequiredImports() *imports.Imports {
	return ti.requiredImports
}

func (ti typeInfo) GoName() string {
	return ti.gti.FullName
}

func (ti typeInfo) GoEffectiveName() string {
	if ti.gti.IsArbitraryType {
		return "interface{}"
	}
	return ti.gti.FullName

}

func (ti typeInfo) GoNameDoc(e Expr) string {
	return ti.gti.NameDoc(e)
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

func (ti typeInfo) GoTypeInfo() *gtypes.Info {
	return ti.gti
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

func (ti typeInfo) TypeMappingsName() string {
	if !ti.IsExported() {
		return ""
	}
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

func (ti typeInfo) PromoteType() string {
	return ti.jti.PromoteType
}

var allTypesSorted = []TypeInfo{}

// This establishes the order in which types are matched by 'case' statements in the "big switch" in goswitch.go. Once established,
// new types cannot be discovered/added.
func SortAllTypes() {
	if len(allTypesSorted) > 0 {
		panic("Attempt to sort all types type after having already sorted all types!!")
	}
	for _, ti := range typesByClojureName {
		if ti.GoName() == "[][]*crypto/x509.Certificate" {
			fmt.Printf("types.go/SortAllTypes: %s == %+v %+v\n", ti.ClojureName(), ti.GoTypeInfo(), ti.ClojureTypeInfo())
		}
		t := ti.GoTypeInfo()
		if t.IsExported && !t.IsArbitraryType {
			allTypesSorted = append(allTypesSorted, ti.(*typeInfo))
		}
	}
	sort.SliceStable(allTypesSorted, func(i, j int) bool {
		i_gti := allTypesSorted[i].GoTypeInfo()
		j_gti := allTypesSorted[j].GoTypeInfo()
		if iSpecificity, jSpecificity := i_gti.Specificity, j_gti.Specificity; iSpecificity != jSpecificity {
			return iSpecificity > jSpecificity
		}
		return i_gti.FullName < j_gti.FullName
	})
}

func AllTypesSorted() []TypeInfo {
	return allTypesSorted
}

func typeKeyForSort(ti TypeInfo) string {
	// Put the pattern (e.g. "%s", "*%s") at the end, to put related types together in a reasonable order.
	return genutils.CombineGoName(ti.GoPackage(), ti.GoBaseName()+ti.GoPattern())
}

func SortedTypeDefinitions(m map[TypeInfo]struct{}, f func(ti TypeInfo)) {
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
