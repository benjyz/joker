package main

import (
	"fmt"
	"github.com/candid82/joker/tools/gostd/godb"
	"github.com/candid82/joker/tools/gostd/gtypes"
	"github.com/candid82/joker/tools/gostd/imports"
	"github.com/candid82/joker/tools/gostd/jtypes"
	. "go/ast"
	"go/token"
	"sort"
	"strings"
)

type TypeInfo interface {
	ArgExtractFunc() string     // Call Extract<this>() for arg with my type
	ArgClojureArgType() string  // Clojure argument type for a Go function arg with my type
	ConvertFromClojure() string // Pattern to convert a (scalar) %s to this type
	ConvertToClojure() string   // Pattern to convert this type to an appropriate Clojure object
	AsJokerObject() string      // Pattern to convert this type to a normal Joker type, or empty string to simply wrap in a GoObject
	JokerName() string
	JokerNameDoc() string
	GoDecl() string
	GoDeclDoc() string
	GoPackage() string
	GoPattern() string
	GoName() string
	GoCode() string
	GoTypeInfo() *gtypes.Info
	TypeSpec() *TypeSpec // Definition, if any, of named type
	GoFile() *godb.GoFile
	DefPos() token.Pos
	Specificity() uint // ConcreteType, else # of methods defined for interface{} (abstract) type
	Ord() uint         // Slot in []*GoTypeInfo and position of case statement in big switch in goswitch.go
	SetOrd(uint)
	TypeMappingsName() string
	Doc() string
	IsNullable() bool // Can an instance of the type == nil (e.g. 'error' type)?
	IsExported() bool
}

const ConcreteType = gtypes.Concrete

type typeInfo struct {
	jti *jtypes.Info
	gti *gtypes.Info
}

type GoTypeInfo struct {
	LocalName                 string       // empty (not a declared type) or the basic type name ("foo" for "x/y.foo")
	FullGoName                string       // empty ("struct {...}" etc.), localName (built-in), path/to/pkg.LocalName, or ABEND if unsupported
	SourceFile                *godb.GoFile // location of the type defintion
	Td                        *TypeSpec
	Type                      TypeInfo // Primary type in the new package
	Where                     token.Pos
	UnderlyingType            Expr             // nil if not a declared type
	ArgClojureType            string           // Can convert this type to a Go function arg with my type
	ArgFromClojureObject      string           // Append this to Clojure object to extract value of my type
	ArgClojureArgType         string           // Clojure argument type for a Go function arg with my type
	ArgExtractFunc            string           // Call Extract<this>() for arg with my type
	ConvertFromClojure        string           // Pattern to convert a (scalar) %s to this type
	ConvertFromClojureImports []imports.Import // Imports needed to support the above
	ConvertFromMap            string           // Pattern to convert a map %s key %s to this type
	ConvertToClojure          string           // Pattern to convert this type to an appropriate Clojure object
	PromoteType               string           // Pattern to convert type to next larger Go type that Joker supports
	ClojureCode               string
	GoCode                    string
	RequiredImports           *imports.Imports
	Uncompleted               bool // Has this type's info been filled in beyond the registration step?
	Custom                    bool // Is this not a builtin Go type?
	Exported                  bool // Is this an exported type?
	Unsupported               bool // Is this unsupported?
	Constructs                bool // Does the convertion from Clojure actually construct (via &sometype{}), returning ptr?
	Nullable                  bool // Can an instance of the type == nil (e.g. 'error' type)?
}

type GoTypeMap map[string]*GoTypeInfo

/* These map fullGoNames to type info. */
var GoTypes = GoTypeMap{}

var typeMap = map[*gtypes.Info]TypeInfo{}

var TypeDefsToGoTypes = map[TypeInfo]*GoTypeInfo{}

func RegisterTypeDecl(ts *TypeSpec, gf *godb.GoFile, pkg string, parentDoc *CommentGroup) bool {
	name := ts.Name.Name
	goTypeName := pkg + "." + name

	// if pkg == "unsafe" && name == "ArbitraryType" {
	// 	if godb.Verbose {
	// 		fmt.Printf("Excluding mythical type %s.%s\n", pkg, name)
	// 	}
	// 	return false
	// }

	if WalkDump {
		fmt.Printf("Type %s at %s:\n", goTypeName, godb.WhereAt(ts.Pos()))
		Print(godb.Fset, ts)
	}

	gtiVec := gtypes.TypeDefine(ts, gf, parentDoc)

	prefix := godb.ClojureNamespaceForPos(godb.Fset.Position(ts.Name.NamePos)) + "/"

	for _, gti := range gtiVec {

		var gt *GoTypeInfo

		gt = registerTypeGOT(gf, goTypeName, ts)
		gt.Td = ts
		gt.Where = ts.Pos()
		gt.RequiredImports = &imports.Imports{}

		jokerName := fmt.Sprintf(gti.Pattern, prefix+name)

		ti := &typeInfo{
			jti: &jtypes.Info{
				ArgClojureArgType:  gt.ArgClojureArgType,
				ConvertFromClojure: gt.ConvertFromClojure,
				ConvertToClojure:   gt.ConvertToClojure,
				JokerName:          jokerName,
				JokerNameDoc:       jokerName,
				AsJokerObject:      gt.ConvertToClojure,
			},
			gti: gti,
		}

		typeMap[gti] = ti

		gt.Type = ti
		TypeDefsToGoTypes[ti] = gt

		if IsExported(name) {
			NumTypes++
			if ti.Specificity() == ConcreteType {
				NumCtableTypes++
			}
		}

		if !strings.Contains(gti.FullName, "[") { // exclude arrays
			ClojureCode[pkg].InitTypes[ti] = struct{}{}
			GoCode[pkg].InitTypes[ti] = struct{}{}
		}
	}

	return true
}

// Maps type-defining Expr to exactly one struct describing that type
var typesByExpr = map[Expr]TypeInfo{}

func BadInfo(err string) typeInfo {
	return typeInfo{
		jti: &jtypes.Info{
			ArgExtractFunc:     err,
			ArgClojureArgType:  err,
			ConvertFromClojure: err + "%0s%0s",
			ConvertToClojure:   err + "%0s%0s",
			AsJokerObject:      err + "%0s%0s",
		},
	}
}

func TypeInfoForExpr(e Expr) TypeInfo {
	gti := gtypes.TypeInfoForExpr(e)
	jti := jtypes.TypeInfoForExpr(e)

	ti := &typeInfo{
		gti: gti,
		jti: jti,
	}

	typeMap[gti] = ti

	return ti
}

func SortedTypeInfoMap(m map[string]*GoTypeInfo, f func(k string, v *GoTypeInfo)) {
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f(k, m[k])
	}
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

func (ti typeInfo) ConvertToClojure() string {
	return ti.jti.ConvertToClojure
}

func (ti typeInfo) AsJokerObject() string {
	if ti.jti == nil {
		return ""
	}
	return ti.jti.AsJokerObject
}

func (ti typeInfo) JokerName() string {
	if ti.jti == nil {
		return ""
	}
	return ti.jti.JokerName
}

func (ti typeInfo) JokerNameDoc() string {
	if ti.jti == nil {
		return ""
	}
	return ti.jti.JokerNameDoc
}

func (ti typeInfo) GoDecl() string {
	return ti.gti.FullName
}

func (ti typeInfo) GoDeclDoc() string {
	return ti.gti.FullName
}

func (ti typeInfo) GoPackage() string {
	return ti.gti.Package
}

func (ti typeInfo) GoPattern() string {
	return ti.gti.Pattern
}

func (ti typeInfo) GoName() string {
	return ti.gti.FullName
}

func (ti typeInfo) GoCode() string {
	return "" // TODO: Probably need something here in some cases? Seems too generic of a name though.
}

func (ti typeInfo) GoTypeInfo() *gtypes.Info { // TODO: Remove when gotypes.go is gone?
	return ti.gti
}

func (ti typeInfo) TypeSpec() *TypeSpec {
	return ti.gti.TypeSpec
}

func (ti typeInfo) GoFile() *godb.GoFile {
	return ti.gti.GoFile
}

func (ti typeInfo) DefPos() token.Pos {
	return ti.gti.DefPos
}

func (ti typeInfo) Specificity() uint {
	return ti.gti.Specificity
}

func (ti typeInfo) Ord() uint {
	return ti.gti.Ord
}

func (ti typeInfo) SetOrd(o uint) {
	ti.gti.Ord = o
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

var allTypesSorted = []TypeInfo{}

// This establishes the order in which types are matched by 'case' statements in the "big switch" in goswitch.go. Once established,
// new types cannot be discovered/added.
func SortAllTypes() {
	if len(allTypesSorted) > 0 {
		panic("Attempt to sort all types type after having already sorted all types!!")
	}
	for _, ti := range typesByExpr {
		t := ti.GoTypeInfo()
		if t.IsExported && (t.Package != "unsafe" || t.LocalName != "ArbitraryType") {
			allTypesSorted = append(allTypesSorted, ti.(*typeInfo))
		}
	}
	sort.SliceStable(allTypesSorted, func(i, j int) bool {
		i_gti := allTypesSorted[i].GoTypeInfo()
		j_gti := allTypesSorted[j].GoTypeInfo()
		if i_gti.Specificity != j_gti.Specificity {
			return i_gti.Specificity > j_gti.Specificity
		}
		return i_gti.FullName < j_gti.FullName
	})
	for ord, t := range allTypesSorted {
		t.SetOrd((uint)(ord))
		ord++
	}
}

func AllTypesSorted() []TypeInfo {
	return allTypesSorted
}

func typeKeyForSort(k string) string {
	if strings.HasPrefix(k, "*") {
		return k[1:] + "*"
	}
	if strings.HasPrefix(k, "[]") {
		return k[2:] + "[]"
	}
	return k
}

func SortedTypeDefinitions(m map[TypeInfo]struct{}, f func(ti TypeInfo)) {
	var keys []string
	vals := map[string]TypeInfo{}
	for k, _ := range m {
		if k != nil {
			key := k.GoTypeInfo().FullName
			keys = append(keys, key)
			vals[key] = k
		}
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return typeKeyForSort(keys[i]) < typeKeyForSort(keys[j])
	})
	for _, k := range keys {
		f(vals[k])
	}
}

func (ti typeInfo) TypeMappingsName() string {
	if !ti.IsExported() {
		return ""
	}
	if ugt := ti.gti.UnderlyingGoType; ugt != nil {
		return "info_PtrTo_" + fmt.Sprintf(ugt.Pattern, ugt.LocalName)
	}
	return "info_" + fmt.Sprintf(ti.GoPattern(), ti.GoName())
}

func init() {
	typeMap[gtypes.Bool] = typeInfo{
		jti: jtypes.Bool,
		gti: gtypes.Bool,
	}
	typeMap[gtypes.Byte] = typeInfo{
		jti: jtypes.Byte,
		gti: gtypes.Byte,
	}
	typeMap[gtypes.Complex128] = typeInfo{
		jti: jtypes.Complex128,
		gti: gtypes.Complex128,
	}
	typeMap[gtypes.Error] = typeInfo{
		jti: jtypes.Error,
		gti: gtypes.Error,
	}
	typeMap[gtypes.Float32] = typeInfo{
		jti: jtypes.Float32,
		gti: gtypes.Float32,
	}
	typeMap[gtypes.Float64] = typeInfo{
		jti: jtypes.Float64,
		gti: gtypes.Float64,
	}
	typeMap[gtypes.Int] = typeInfo{
		jti: jtypes.Int,
		gti: gtypes.Int,
	}
	typeMap[gtypes.Int32] = typeInfo{
		jti: jtypes.Int32,
		gti: gtypes.Int32,
	}
	typeMap[gtypes.Int64] = typeInfo{
		jti: jtypes.Int64,
		gti: gtypes.Int64,
	}
	typeMap[gtypes.Rune] = typeInfo{
		jti: jtypes.Rune,
		gti: gtypes.Rune,
	}
	typeMap[gtypes.String] = typeInfo{
		jti: jtypes.String,
		gti: gtypes.String,
	}
	typeMap[gtypes.UInt] = typeInfo{
		jti: jtypes.UInt,
		gti: gtypes.UInt,
	}
	typeMap[gtypes.UInt16] = typeInfo{
		jti: jtypes.UInt16,
		gti: gtypes.UInt16,
	}
	typeMap[gtypes.UInt32] = typeInfo{
		jti: jtypes.UInt32,
		gti: gtypes.UInt32,
	}
	typeMap[gtypes.UInt64] = typeInfo{
		jti: jtypes.UInt64,
		gti: gtypes.UInt64,
	}
	typeMap[gtypes.UInt8] = typeInfo{
		jti: jtypes.UInt8,
		gti: gtypes.UInt8,
	}
	typeMap[gtypes.UIntPtr] = typeInfo{
		jti: jtypes.UIntPtr,
		gti: gtypes.UIntPtr,
	}
}
