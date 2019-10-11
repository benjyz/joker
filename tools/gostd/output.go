package main

import (
	"bufio"
	"fmt"
	. "github.com/candid82/joker/tools/gostd/gowalk"
	"github.com/candid82/joker/tools/gostd/imports"
	. "github.com/candid82/joker/tools/gostd/types"
	. "github.com/candid82/joker/tools/gostd/utils"
	"go/doc"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var currentTimeAndVersion = ""
var noTimeAndVersion = false

func curTimeAndVersion() string {
	if noTimeAndVersion {
		return "(omitted for testing)"
	}
	if currentTimeAndVersion == "" {
		by, _ := time.Now().MarshalText()
		currentTimeAndVersion = string(by) + " by version " + VERSION
	}
	return currentTimeAndVersion
}

func RegisterPackages(pkgs []string, jokerSourceDir string) {
	updateCustomLibsGo(pkgs, filepath.Join(jokerSourceDir, "custom.go"))
}

func RegisterJokerFiles(jokerFiles []string, jokerSourceDir string) {
	updateCustomLibsJoker(jokerFiles, filepath.Join(jokerSourceDir, "core", "data", "customlibs.joke"))
}

func RegisterGoTypeSwitch(types []*TypeDefInfo, jokerSourceDir string, outputCode bool) {
	updateGoTypeSwitch(types, filepath.Join(jokerSourceDir, "core", "goswitch.go"), outputCode)
}

// E.g.: \t_ "github.com/candid82/joker/std/go/std/net"
func updateCustomLibsGo(pkgs []string, f string) {
	if Verbose {
		fmt.Printf("Adding %d custom imports to %s\n", len(pkgs), filepath.ToSlash(f))
	}

	var m string
	if len(pkgs) > 0 {
		m = "// Auto-modified by gostd at " + curTimeAndVersion()
	} else {
		m = "// Placeholder for custom libraries. Overwritten by gostd."
	}

	m += `

package main
`

	if len(pkgs) > 0 {
		newImports := `

import (
`
		importPrefix := "\t_ \"github.com/candid82/joker/std/go/std/"
		for _, p := range pkgs {
			newImports += importPrefix + p + "\"\n"
		}
		newImports += `)
`
		m += newImports
	}

	err := ioutil.WriteFile(f, []byte(m), 0777)
	Check(err)
}

func updateCustomLibsJoker(pkgs []string, f string) {
	if Verbose {
		fmt.Printf("Adding %d custom loaded libraries to %s\n", len(pkgs), filepath.ToSlash(f))
	}

	var m string
	if len(pkgs) > 0 {
		m = ";; Auto-modified by gostd at " + curTimeAndVersion()
	} else {
		m = ";; Placeholder for custom libraries. Overwritten by gostd."
	}

	m += `

(def ^:dynamic
  ^{:private true
    :doc "A set of symbols representing loaded custom libs"}
  *custom-libs* #{
`

	const importPrefix = " 'go.std."
	for _, p := range pkgs {
		m += "    " + importPrefix + strings.Replace(p, "/", ".", -1) + "\n"
	}
	m += `    })
`

	err := ioutil.WriteFile(f, []byte(m), 0777)
	Check(err)
}

func updateGoTypeSwitch(types []*TypeDefInfo, f string, outputCode bool) {
	if Verbose {
		fmt.Printf("Adding %d types to %s\n", len(types), filepath.ToSlash(f))
	}

	var pattern string
	if len(types) > 0 {
		pattern = "// Auto-modified by gostd at " + curTimeAndVersion()
	} else {
		pattern = "// Placeholder for big Go switch on types. Overwritten by gostd."
	}

	pattern += `

package core

import (
%s)

var GoTypesVec [%d]*GoTypeInfo

func SwitchGoType(g interface{}) *GoTypeInfo {
	switch g.(type) {
%s	}
	return nil
}
`

	var cases string
	for _, t := range types {
		pkgPlusSeparator := ""
		if t.GoPackage != "" {
			pkgPlusSeparator = t.GoPackage + "."
		}
		cases += fmt.Sprintf("\tcase %s%s%s:\n\t\treturn GoTypesVec[%d]\n", t.GoPrefix, pkgPlusSeparator, t.GoName, t.Ord)
	}

	m := fmt.Sprintf(pattern, "", len(AllSorted()), cases)

	err := ioutil.WriteFile(f, []byte(m), 0777)
	// Ignore error if outputting code to stdout:
	if !outputCode {
		Check(err)
	}

	if outputCode {
		fmt.Println("Generated file goswitch.go:")
		fmt.Print(m)
	}
}

func packageQuotedImportList(pi imports.Imports, prefix string) string {
	imports := ""
	SortedPackageImports(pi,
		func(k, local, full string) {
			if local == "" {
				imports += prefix + `"` + k + `"`
			} else {
				imports += prefix + local + ` "` + k + `"`
			}
		})
	return imports
}

func outputClojureCode(pkgDirUnix string, v CodeInfo, jokerLibDir string, outputCode, generateEmpty bool) {
	var out *bufio.Writer
	var unbuf_out *os.File

	if jokerLibDir != "" && jokerLibDir != "-" &&
		(generateEmpty || PackagesInfo[pkgDirUnix].NonEmpty) {
		jf := filepath.Join(jokerLibDir, filepath.FromSlash(pkgDirUnix)+".joke")
		var e error
		e = os.MkdirAll(filepath.Dir(jf), 0777)
		unbuf_out, e = os.Create(jf)
		Check(e)
	} else if generateEmpty || PackagesInfo[pkgDirUnix].NonEmpty {
		unbuf_out = os.Stdout
	}
	if unbuf_out != nil {
		out = bufio.NewWriterSize(unbuf_out, 16384)
	}

	pi := PackagesInfo[pkgDirUnix]

	if out != nil {
		importPath, _ := filepath.Abs("/")
		myDoc := doc.New(pi.Pkg, importPath, doc.AllDecls)
		pkgDoc := fmt.Sprintf("Provides a low-level interface to the %s package.", pkgDirUnix)
		if myDoc.Doc != "" {
			pkgDoc += "\n\n" + myDoc.Doc
		}

		fmt.Fprintf(out,
			`;;;; Auto-generated by gostd at `+curTimeAndVersion()+`, do not edit!!

(ns
  ^{:go-imports [%s]
    :doc %s
    :empty %s}
  %s)
`,
			strings.TrimPrefix(packageQuotedImportList(*pi.ImportsAutoGen, " "), " "),
			strconv.Quote(pkgDoc),
			func() string {
				if pi.NonEmpty {
					return "false"
				} else {
					return "true"
				}
			}(),
			"go.std."+strings.Replace(pkgDirUnix, "/", ".", -1))
	}

	SortedConstantInfoMap(v.Constants,
		func(c string, ci *ConstantInfo) {
			if outputCode {
				fmt.Printf("JOKER CONSTANT %s from %s:%s\n", c, ci.SourceFile.Name, ci.Def)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ci.Def)
			}
		})

	SortedVariableInfoMap(v.Variables,
		func(c string, ci *VariableInfo) {
			if outputCode {
				fmt.Printf("JOKER VARIABLE %s from %s:%s\n", c, ci.SourceFile.Name, ci.Def)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ci.Def)
			}
		})

	SortedTypeInfoMap(v.Types,
		func(t string, ti *GoTypeInfo) {
			if outputCode {
				fmt.Printf("JOKER TYPE %s from %s:%s\n", t, ti.SourceFile.Name, ti.ClojureCode)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ti.ClojureCode)
			}
		})

	SortedCodeMap(v,
		func(f string, w *FnCodeInfo) {
			if outputCode {
				fmt.Printf("JOKER FUNC %s.%s from %s:%s\n",
					pkgDirUnix, f, w.SourceFile.Name, w.FnCode)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(w.FnCode)
			}
		})

	SortedTypeDefinitions(v.InitTypes,
		func(tdi *TypeDefInfo) {
			tmn := tdi.TypeMappingsName()
			if tmn == "" || tdi.LocalName == "" || tdi.IsPrivate {
				return
			}
			typeDoc := tdi.Doc
			fnCode := fmt.Sprintf(`
(def
  ^{:doc %s
    :added "1.0"
    :tag "GoType"
    :go "&%s"}
  %s)
`,
				strconv.Quote(typeDoc), tmn, tdi.LocalName)
			if outputCode {
				fmt.Printf("JOKER TYPE %s:%s\n",
					tdi.FullName, fnCode)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(fnCode)
			}
		})

	if out != nil {
		out.Flush()
		if unbuf_out != os.Stdout {
			unbuf_out.Close()
		}
	}
}

func outputGoCode(pkgDirUnix string, v CodeInfo, jokerLibDir string, outputCode, generateEmpty bool) {
	pkgBaseName := path.Base(pkgDirUnix)
	pi := PackagesInfo[pkgDirUnix]
	PackagesInfo[pkgDirUnix].HasGoFiles = true
	pkgDirNative := filepath.FromSlash(pkgDirUnix)

	var out *bufio.Writer
	var unbuf_out *os.File

	if jokerLibDir != "" && jokerLibDir != "-" &&
		(generateEmpty || PackagesInfo[pkgDirUnix].NonEmpty) {
		gf := filepath.Join(jokerLibDir, pkgDirNative,
			pkgBaseName+"_native.go")
		var e error
		e = os.MkdirAll(filepath.Dir(gf), 0777)
		Check(e)
		unbuf_out, e = os.Create(gf)
		Check(e)
	} else if generateEmpty || PackagesInfo[pkgDirUnix].NonEmpty {
		unbuf_out = os.Stdout
	}
	if unbuf_out != nil {
		out = bufio.NewWriterSize(unbuf_out, 16384)
	}

	if out != nil {
		fmt.Fprintf(out,
			`// Auto-generated by gostd at `+curTimeAndVersion()+`, do not edit!!

package %s

import (%s
)
`,
			pkgBaseName,
			packageQuotedImportList(*pi.ImportsNative, "\n\t"))
	}

	SortedTypeInfoMap(v.Types,
		func(t string, ti *GoTypeInfo) {
			if outputCode {
				fmt.Printf("GO TYPE %s from %s:%s\n", t, ti.SourceFile.Name, ti.GoCode)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ti.GoCode)
			}
		})

	SortedCodeMap(v,
		func(f string, w *FnCodeInfo) {
			if outputCode {
				fmt.Printf("GO FUNC %s.%s from %s:%s\n",
					pkgDirUnix, f, w.SourceFile.Name, w.FnCode)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(w.FnCode)
			}
		})

	SortedTypeDefinitions(v.InitTypes,
		func(tdi *TypeDefInfo) {
			tmn := tdi.TypeMappingsName()
			if tmn == "" || tdi.IsPrivate {
				return
			}
			tmn = fmt.Sprintf("var %s GoTypeInfo\n", tmn)
			if outputCode && tmn != "" {
				fmt.Printf("GO VARDEF FOR TYPE %s from %s:\n%s\n", tdi.FullName, WhereAt(tdi.DefPos), tmn)
			}
			if out != nil && unbuf_out != os.Stdout && tmn != "" {
				out.WriteString(tmn)
			}
		})

	const initInfoTemplate = `
	%s = GoTypeInfo{Name: "%s",
		GoType: &GoType{T: &%s},
		Members: GoMembers{
%s		}}

`

	if out != nil {
		out.WriteString("\nfunc initNative() {\n")
	}
	SortedTypeDefinitions(v.InitTypes,
		func(tdi *TypeDefInfo) {
			tmn := tdi.TypeMappingsName()
			if tmn == "" || tdi.IsPrivate {
				return
			}
			k1 := tdi.FullName
			mem := ""
			SortedFnCodeInfo(v.InitVars[tdi], // Will always be populated
				func(c string, r *FnCodeInfo) {
					doc := r.FnDoc
					g := r.FnCode
					mem += fmt.Sprintf(`
			"%s": MakeGoReceiver("%s", %s, %s, %s, NewVectorFrom(%s)),
`[1:],
						c, c, g, strconv.Quote(CommentGroupAsString(doc)), strconv.Quote("1.0"), paramsAsSymbolVec(r.FnDecl.Type.Params))
				})
			o := fmt.Sprintf(initInfoTemplate[1:], tmn, k1, tmn, mem)
			if outputCode {
				fmt.Printf("GO INFO FOR TYPE %s from %s:\n%s\n", tdi.FullName, WhereAt(tdi.DefPos), o)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(o)
			}
		})

	SortedTypeDefinitions(v.InitTypes,
		func(tdi *TypeDefInfo) {
			reflectPackageImport, reflectPattern := tdi.TypeReflected()
			tmn := tdi.TypeMappingsName()
			if reflectPattern == "" || tmn == "" || tdi.IsPrivate {
				return
			}
			reflectLocal := imports.AddImport(PackagesInfo[pkgDirUnix].ImportsNative, "", reflectPackageImport, true)
			o := fmt.Sprintf("\tGoTypes[%s] = &%s\n", fmt.Sprintf(reflectPattern, reflectLocal), tmn)
			if outputCode {
				fmt.Printf("GO VECSET FOR TYPE %s from %s:\n%s\n", tdi.FullName, WhereAt(tdi.DefPos), o)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(o)
			}
		})

	if out != nil {
		out.WriteString("}\n")
		if unbuf_out == os.Stdout {
			out.WriteString("\n") // separate from next "file" output for testing
		}
	}

	if out != nil {
		out.Flush()
		if unbuf_out != os.Stdout {
			unbuf_out.Close()
		}
	}
}

func OutputPackageCode(jokerLibDir string, outputCode, generateEmpty bool) {
	SortedPackageMap(ClojureCode,
		func(pkgDirUnix string, v CodeInfo) {
			outputClojureCode(pkgDirUnix, v, jokerLibDir, outputCode, generateEmpty)
		})

	SortedPackageMap(GoCode,
		func(pkgDirUnix string, v CodeInfo) {
			outputGoCode(pkgDirUnix, v, jokerLibDir, outputCode, generateEmpty)
		})
}
