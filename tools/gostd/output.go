package main

import (
	"bufio"
	"fmt"
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

func registerPackages(pkgs []string, jokerSourceDir string) {
	updateCustomLibsGo(pkgs, filepath.Join(jokerSourceDir, "custom.go"))
}

func registerJokerFiles(jokerFiles []string, jokerSourceDir string) {
	updateCustomLibsJoker(jokerFiles, filepath.Join(jokerSourceDir, "core", "data", "customlibs.joke"))
}

// E.g.: \t_ "github.com/candid82/joker/std/go/std/net"
func updateCustomLibsGo(pkgs []string, f string) {
	if verbose {
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
	check(err)
}

func updateCustomLibsJoker(pkgs []string, f string) {
	if verbose {
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
	check(err)
}

func packageQuotedImportList(pi packageImports, prefix string) string {
	imports := ""
	sortedPackageImports(pi,
		func(k, local, full string) {
			if local == "" {
				imports += prefix + `"` + k + `"`
			} else {
				imports += prefix + local + ` "` + k + `"`
			}
		})
	return imports
}

func outputClojureCode(pkgDirUnix string, v codeInfo, jokerLibDir string, outputCode, generateEmpty bool) {
	var out *bufio.Writer
	var unbuf_out *os.File

	if jokerLibDir != "" && jokerLibDir != "-" &&
		(generateEmpty || packagesInfo[pkgDirUnix].nonEmpty) {
		jf := filepath.Join(jokerLibDir, filepath.FromSlash(pkgDirUnix)+".joke")
		var e error
		e = os.MkdirAll(filepath.Dir(jf), 0777)
		unbuf_out, e = os.Create(jf)
		check(e)
	} else if generateEmpty || packagesInfo[pkgDirUnix].nonEmpty {
		unbuf_out = os.Stdout
	}
	if unbuf_out != nil {
		out = bufio.NewWriterSize(unbuf_out, 16384)
	}

	pi := packagesInfo[pkgDirUnix]

	if out != nil {
		importPath, _ := filepath.Abs("/")
		myDoc := doc.New(pi.pkg, importPath, doc.AllDecls)
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
			strings.TrimPrefix(packageQuotedImportList(*pi.importsAutoGen, " "), " "),
			strconv.Quote(pkgDoc),
			func() string {
				if pi.nonEmpty {
					return "false"
				} else {
					return "true"
				}
			}(),
			"go.std."+strings.Replace(pkgDirUnix, "/", ".", -1))
	}

	sortedConstantInfoMap(v.constants,
		func(c string, ci *constantInfo) {
			if outputCode {
				fmt.Printf("JOKER CONSTANT %s from %s:%s\n", c, ci.sourceFile.name, ci.def)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ci.def)
			}
		})

	sortedVariableInfoMap(v.variables,
		func(c string, ci *variableInfo) {
			if outputCode {
				fmt.Printf("JOKER VARIABLE %s from %s:%s\n", c, ci.sourceFile.name, ci.def)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ci.def)
			}
		})

	sortedTypeInfoMap(v.types,
		func(t string, ti *goTypeInfo) {
			if outputCode {
				fmt.Printf("JOKER TYPE %s from %s:%s\n", t, ti.sourceFile.name, ti.clojureCode)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ti.clojureCode)
			}
		})

	sortedCodeMap(v,
		func(f string, w fnCodeInfo) {
			if outputCode {
				fmt.Printf("JOKER FUNC %s.%s from %s:%s\n",
					pkgDirUnix, f, w.sourceFile.name, w.fnCode)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(w.fnCode)
			}
		})

	if out != nil {
		out.Flush()
		if unbuf_out != os.Stdout {
			unbuf_out.Close()
		}
	}
}

func outputGoCode(pkgDirUnix string, v codeInfo, jokerLibDir string, outputCode, generateEmpty bool) {
	pkgBaseName := path.Base(pkgDirUnix)
	pi := packagesInfo[pkgDirUnix]
	packagesInfo[pkgDirUnix].hasGoFiles = true
	pkgDirNative := filepath.FromSlash(pkgDirUnix)

	var out *bufio.Writer
	var unbuf_out *os.File

	if jokerLibDir != "" && jokerLibDir != "-" &&
		(generateEmpty || packagesInfo[pkgDirUnix].nonEmpty) {
		gf := filepath.Join(jokerLibDir, pkgDirNative,
			pkgBaseName+"_native.go")
		var e error
		e = os.MkdirAll(filepath.Dir(gf), 0777)
		check(e)
		unbuf_out, e = os.Create(gf)
		check(e)
	} else if generateEmpty || packagesInfo[pkgDirUnix].nonEmpty {
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
			packageQuotedImportList(*pi.importsNative, "\n\t"))
	}

	sortedTypeInfoMap(v.types,
		func(t string, ti *goTypeInfo) {
			if outputCode {
				fmt.Printf("GO TYPE %s from %s:%s\n", t, ti.sourceFile.name, ti.goCode)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(ti.goCode)
			}
		})

	sortedCodeMap(v,
		func(f string, w fnCodeInfo) {
			if outputCode {
				fmt.Printf("GO FUNC %s.%s from %s:%s\n",
					pkgDirUnix, f, w.sourceFile.name, w.fnCode)
			}
			if out != nil && unbuf_out != os.Stdout {
				out.WriteString(w.fnCode)
			}
		})

	sortedStringMap(v.initTypes,
		func(k1, k2 string) {
			out.WriteString(fmt.Sprintf("var %s GoTypeInfo\n", k2))
		})

	const initInfoTemplate = `
	%s = GoTypeInfo{Name: "%s",
		GoType: GoType{T: &%s},
		Members: GoMembers{
%s		}}

`

	if out != nil {
		out.WriteString("\nfunc initNative() {\n")
	}
	sortedStringMap(v.initTypes,
		func(k1, k2 string) {
			mem := ""
			sortedStringMap(v.initVars[k2], // Will always be populated
				func(c, g string) {
					mem += fmt.Sprintf(`
			"%s": MakeGoReceiver("%s", %s),
`[1:],
						c, c, g)
				})
			out.WriteString(fmt.Sprintf(initInfoTemplate[1:], k2, v.initTypesFullName[k1], k2, mem))
		})

	const internTypeTemplate = `
        %sNamespace.InternVar("%s", MakeGoType(&info_%s),
                MakeMeta(
                        nil,
                        %s, "%s"))

`

	/*
		sortedStringMap(v.initTypes,
			func(k, v string) {
				out.WriteString(fmt.Sprintf(internTypeTemplate[1:], "net", "MX", "MX", "doc for MX", "1.0"))
			})
	*/
	sortedStringMap(v.initTypes,
		func(k, v string) {
			out.WriteString(fmt.Sprintf("\tGoTypes[%s] = &%s\n", k, v))
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

func outputPackageCode(jokerLibDir string, outputCode, generateEmpty bool) {
	sortedPackageMap(clojureCode,
		func(pkgDirUnix string, v codeInfo) {
			outputClojureCode(pkgDirUnix, v, jokerLibDir, outputCode, generateEmpty)
		})

	sortedPackageMap(goCode,
		func(pkgDirUnix string, v codeInfo) {
			outputGoCode(pkgDirUnix, v, jokerLibDir, outputCode, generateEmpty)
		})
}
