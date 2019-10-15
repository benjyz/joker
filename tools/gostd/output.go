package main

import (
	"bufio"
	"fmt"
	. "github.com/candid82/joker/tools/gostd/godb"
	"github.com/candid82/joker/tools/gostd/gowalk"
	"github.com/candid82/joker/tools/gostd/imports"
	. "github.com/candid82/joker/tools/gostd/types"
	. "github.com/candid82/joker/tools/gostd/utils"
	"go/doc"
	"go/token"
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
	if gowalk.Verbose {
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
	if gowalk.Verbose {
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
		m += "    " + importPrefix + strings.ReplaceAll(p, "/", ".") + "\n"
	}
	m += `    })
`

	err := ioutil.WriteFile(f, []byte(m), 0777)
	Check(err)
}

func updateGoTypeSwitch(types []*TypeDefInfo, f string, outputCode bool) {
	if gowalk.Verbose {
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

import (%s
)

var GoTypesVec [%d]*GoTypeInfo

func SwitchGoType(g interface{}) int {
	switch g.(type) {
%s	}
	return -1
}
`

	var cases string
	var importeds = &imports.Imports{}
	for _, t := range types {
		pkgPlusSeparator := ""
		if t.GoPackage != "" {
			pkgPlusSeparator = imports.AddImport(importeds, "", t.GoPackage, true) + "."
		}
		specificity := ""
		if t.Specificity != Concrete {
			specificity = fmt.Sprintf("  // Specificity=%d", t.Specificity)
		}
		cases += fmt.Sprintf("\tcase %s%s%s:%s\n\t\treturn %d\n", t.GoPrefix, pkgPlusSeparator, t.GoName, specificity, t.Ord)
	}

	m := fmt.Sprintf(pattern, imports.QuotedImportList(importeds, "\n\t"), len(types), cases)

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

func outputClojureCode(pkgDirUnix string, v gowalk.CodeInfo, jokerLibDir string, outputCode, generateEmpty bool) {
	var out *bufio.Writer
	var unbuf_out *os.File

	if jokerLibDir != "" && jokerLibDir != "-" &&
		(generateEmpty || gowalk.PackagesInfo[pkgDirUnix].NonEmpty) {
		jf := filepath.Join(jokerLibDir, filepath.FromSlash(pkgDirUnix)+".joke")
		var e error
		e = os.MkdirAll(filepath.Dir(jf), 0777)
		unbuf_out, e = os.Create(jf)
		Check(e)
	} else if generateEmpty || gowalk.PackagesInfo[pkgDirUnix].NonEmpty {
		unbuf_out = os.Stdout
	}
	if unbuf_out != nil {
		out = bufio.NewWriterSize(unbuf_out, 16384)
	}

	pi := gowalk.PackagesInfo[pkgDirUnix]

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
			strings.TrimPrefix(imports.QuotedImportList(pi.ImportsAutoGen, " "), " "),
			strconv.Quote(pkgDoc),
			func() string {
				if pi.NonEmpty {
					return "false"
				} else {
					return "true"
				}
			}(),
			"go.std."+strings.ReplaceAll(pkgDirUnix, "/", "."))
	}

	gowalk.SortedConstantInfoMap(v.Constants,
		func(c string, ci *gowalk.ConstantInfo) {
			if outputCode {
				fmt.Printf("JOKER CONSTANT %s from %s:%s\n", c, ci.SourceFile.Name, ci.Def)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ci.Def)
			}
		})

	gowalk.SortedVariableInfoMap(v.Variables,
		func(c string, ci *gowalk.VariableInfo) {
			if outputCode {
				fmt.Printf("JOKER VARIABLE %s from %s:%s\n", c, ci.SourceFile.Name, ci.Def)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ci.Def)
			}
		})

	gowalk.SortedTypeInfoMap(v.Types,
		func(t string, ti *gowalk.GoTypeInfo) {
			if outputCode {
				fmt.Printf("JOKER TYPE %s from %s:%s\n", t, ti.SourceFile.Name, ti.ClojureCode)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ti.ClojureCode)
			}
		})

	gowalk.SortedCodeMap(v,
		func(f string, w *gowalk.FnCodeInfo) {
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
			if tmn == "" || tdi.LocalName == "" || !tdi.IsExported {
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

func outputGoCode(pkgDirUnix string, v gowalk.CodeInfo, jokerLibDir string, outputCode, generateEmpty bool) {
	pkgBaseName := path.Base(pkgDirUnix)
	pi := gowalk.PackagesInfo[pkgDirUnix]
	pi.HasGoFiles = true
	pkgDirNative := filepath.FromSlash(pkgDirUnix)

	var out *bufio.Writer
	var unbuf_out *os.File

	if jokerLibDir != "" && jokerLibDir != "-" &&
		(generateEmpty || pi.NonEmpty) {
		gf := filepath.Join(jokerLibDir, pkgDirNative,
			pkgBaseName+"_native.go")
		var e error
		e = os.MkdirAll(filepath.Dir(gf), 0777)
		Check(e)
		unbuf_out, e = os.Create(gf)
		Check(e)
	} else if generateEmpty || pi.NonEmpty {
		unbuf_out = os.Stdout
	}
	if unbuf_out != nil {
		out = bufio.NewWriterSize(unbuf_out, 16384)
	}

	// First, figure out what other packages need to be imported,
	// before the import statement is generated.
	ensure := ""
	imports.SortedOriginalPackageImports(pi.Pkg,
		gowalk.LegitimateImport,
		func(imp string, pos token.Pos) {
			ns := ClojureNamespaceForDirname(imp)
			if ns == pi.ClojureNameSpace {
				return // it me
			}

			imports.AddImport(pi.ImportsNative, ".", "github.com/candid82/joker/core", false)

			ensure += fmt.Sprintf("\tEnsureLoaded(\"%s\")  // E.g. from: %s\n", ns, WhereAt(pos))
		})

	if out != nil {
		fmt.Fprintf(out,
			`// Auto-generated by gostd at `+curTimeAndVersion()+`, do not edit!!

package %s

import (%s
)
`,
			pkgBaseName,
			imports.QuotedImportList(pi.ImportsNative, "\n\t"))
	}

	gowalk.SortedTypeInfoMap(v.Types,
		func(t string, ti *gowalk.GoTypeInfo) {
			ctor := ""
			if c, found := Ctors[ti.TypeDefInfo]; found && c[0] != '/' {
				ctor = c
			}
			if outputCode {
				fmt.Printf("GO TYPE %s from %s:%s\n", t, ti.SourceFile.Name, ti.GoCode+ctor)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ti.GoCode + ctor)
			}
		})

	gowalk.SortedCodeMap(v,
		func(f string, w *gowalk.FnCodeInfo) {
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
			if tmn == "" || !tdi.IsExported {
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
%s		Members: GoMembers{
%s		}}

`

	if out != nil {
		out.WriteString("\nfunc initNative() {\n")
	}

	SortedTypeDefinitions(v.InitTypes,
		func(tdi *TypeDefInfo) {
			tmn := tdi.TypeMappingsName()
			if tmn == "" || !tdi.IsExported {
				return
			}
			k1 := tdi.FullName
			ctor := ""
			if c, found := CtorNames[tdi]; found {
				ctor = fmt.Sprintf(`
		Ctor: %s,
`[1:],
					c)
			}
			mem := ""
			gowalk.SortedFnCodeInfo(v.InitVars[tdi], // Will always be populated
				func(c string, r *gowalk.FnCodeInfo) {
					doc := r.FnDoc
					g := r.FnCode
					mem += fmt.Sprintf(`
			"%s": MakeGoReceiver("%s", %s, %s, %s, NewVectorFrom(%s)),
`[1:],
						c, c, g, strconv.Quote(CommentGroupAsString(doc)), strconv.Quote("1.0"), paramsAsSymbolVec(r.Params))
				})
			o := fmt.Sprintf(initInfoTemplate[1:], tmn, k1, tmn, ctor, mem)
			if outputCode {
				fmt.Printf("GO INFO FOR TYPE %s from %s:\n%s\n", tdi.FullName, WhereAt(tdi.DefPos), o)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(o)
			}
		})

	SortedTypeDefinitions(v.InitTypes,
		func(tdi *TypeDefInfo) {
			tmn := tdi.TypeMappingsName()
			if tmn == "" || !tdi.IsExported {
				return
			}
			o := fmt.Sprintf("\tGoTypesVec[%d] = &%s\n", tdi.Ord, tmn)
			if outputCode {
				fmt.Printf("GO VECSET FOR TYPE %s from %s:\n%s\n", tdi.FullName, WhereAt(tdi.DefPos), o)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(o)
			}
		})

	if ensure != "" {
		if outputCode {
			fmt.Printf("GO ENSURE-LOADED FOR %s:\n%s\n", pi.Pkg.Name, ensure)
		}
		if out != nil && unbuf_out != os.Stdout {
			out.WriteString(ensure)
		}
	}

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
	gowalk.SortedPackageMap(gowalk.ClojureCode,
		func(pkgDirUnix string, v gowalk.CodeInfo) {
			outputClojureCode(pkgDirUnix, v, jokerLibDir, outputCode, generateEmpty)
		})

	gowalk.SortedPackageMap(gowalk.GoCode,
		func(pkgDirUnix string, v gowalk.CodeInfo) {
			outputGoCode(pkgDirUnix, v, jokerLibDir, outputCode, generateEmpty)
		})
}
