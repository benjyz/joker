package main

import (
	"fmt"
	. "go/ast"
	"go/token"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func whereAt(p token.Pos) string {
	return fmt.Sprintf("%s", fset.Position(p).String())
}

func fileAt(p token.Pos) string {
	return token.Position{Filename: fset.Position(p).Filename,
		Offset: 0, Line: 0, Column: 0}.String()
}

func unix(p string) string {
	return filepath.ToSlash(p)
}

func commentGroupInQuotes(doc *CommentGroup, jokIn, jokOut, goIn, goOut string) string {
	var d string
	if doc != nil {
		d = doc.Text()
	}
	if goIn != "" {
		if d != "" {
			d = strings.Trim(d, " \t\n") + "\n\n"
		}
		d += "Go input arguments: " + goIn
	}
	if goOut != "" {
		if d != "" {
			d = strings.Trim(d, " \t\n") + "\n\n"
		}
		d += "Go return type: " + goOut
	}
	if jokIn != "" {
		if d != "" {
			d = strings.Trim(d, " \t\n") + "\n\n"
		}
		d += "Joker input arguments: " + jokIn
	}
	if jokOut != "" {
		if d != "" {
			d = strings.Trim(d, " \t\n") + "\n\n"
		}
		d += "Joker return type: " + jokOut
	}
	return `  ` + strings.Trim(strconv.Quote(d), " \t\n") + "\n"
}

func paramNameAsClojure(n string) string {
	return n
}

func paramNameAsGo(p string) string {
	return p
}

func replaceAll(string, from, to string) string {
	return strings.Replace(string, from, to, -1)
}

func fullTypeNameAsClojure(t string) string {
	if t[0] == '_' {
		t = t[1:]
	}
	return "go.std." + replaceAll(replaceAll(replaceAll(t, ".", ":"), "/", "."), ":", "/")
}

// Given an input package name such as "foo/bar" and typename
// "bletch", decides whether to return (for 'code' and 'cl2gol') just
// "_bar.bletch" and "bletch" if the package being compiled will be
// implementing Go's package of the same name (in this case, the
// generated file will be foo/bar_native.go and start with "package
// bar"); or, to return (for both) simply "bar.bletch" and ensure
// "foo/bar" is imported (implicitly as "bar", assuming no
// conflicts). NOTE: As a side effect, updates imports needed by the
// function.
func fullPkgNameAsGoType(fn *funcInfo, fullPkgName, baseTypeName string) (code, doc, cl2gol string) {
	curPkgName := fn.sourceFile.pkgDirUnix
	if curPkgName == fullPkgName {
		code = "_" + path.Base(fullPkgName) + "." + baseTypeName
		doc = baseTypeName
		cl2gol = path.Base(fullPkgName) + "." + baseTypeName
		fn.refersToSelf = true
		return
	}
	doc = path.Base(fullPkgName) + "." + baseTypeName
	code = doc
	cl2gol = "ABEND987(genutils.go: imports not yet supported: " + doc + ")"
	return
}

func isPredeclaredIdentifier(s string) bool {
	switch s {
	case "append", "bool", "byte", "cap", "close", "complex", "complex128", "complex64", "copy", "delete", "error", "false", "float32", "float64", "imag", "int", "int16", "int32", "int64", "int8", "iota", "len", "make", "new", "nil", "panic", "print", "println", "real", "recover", "rune", "string", "true", "uint", "uint16", "uint32", "uint64", "uint8", "uintptr":
		return true
	default:
		return false
	}
}

func funcNameAsGoPrivate(f string) string {
	s := strings.ToLower(f[0:1]) + f[1:]
	if token.Lookup(s).IsKeyword() || isPredeclaredIdentifier(s) {
		s = "_" + s
	}
	return s
}

func isPrivate(p string) bool {
	return !unicode.IsUpper(rune(p[0]))
}

func reverseJoin(a []string, infix string) string {
	j := ""
	for idx := len(a) - 1; idx >= 0; idx-- {
		if idx != len(a)-1 {
			j += infix
		}
		j += a[idx]
	}
	return j
}

var genSymIndex = map[string]int{}

func genSym(pre string) string {
	var idx int
	if i, ok := genSymIndex[pre]; ok {
		idx = i + 1
	} else {
		idx = 1
	}
	genSymIndex[pre] = idx
	return fmt.Sprintf("%s%d", pre, idx)
}

func genSymReset() {
	genSymIndex = map[string]int{}
}

func exprIsUseful(rtn string) bool {
	return rtn != "NIL"
}

// Generates code that, at run time, tests each of the onlyIf's and, if all true, returns the expr; else returns NIL.
func wrapOnlyIfs(onlyIf string, e string) string {
	if len(onlyIf) == 0 {
		return e
	}
	return "func() Object { if " + onlyIf + " { return " + e + " } else { return NIL } }()"
}

// Add one level of indent to each line
func indentedCode(c string) string {
	return "\t" + strings.Replace(c, "\n", "\n\t", -1)
}

func wrapStmtOnlyIfs(indent, v, t, e string, onlyIf string, c string, out *string) string {
	if len(onlyIf) == 0 {
		*out = v
		return indent + v + " := " + e + "\n" + c
	}
	*out = "_obj" + v
	return indent + "var " + *out + " Object\n" +
		indent + "if " + onlyIf + " {\n" +
		indent + "\t" + v + " := " + e + "\n" +
		strings.TrimRight(indentedCode(c), "\t") +
		indent + "\t" + *out + " = Object(" + v + ")\n" +
		indent + "} else {\n" +
		indent + "\t" + *out + " = NIL\n" +
		indent + "}\n"
}

// Return a form of the return type as supported by generate-std.joke,
// or empty string if not supported (which will trigger attempting to
// generate appropriate code for *_native.go). gol either passes
// through or "Object" is returned for it if cl is returned as empty.
func clojureReturnTypeForGenerateCustom(in_cl, in_gol string) (cl, gol string) {
	switch in_cl {
	case "String", "Int", "Byte", "Double", "Bool", "Time", "Error": // TODO: Have tested only String so far
		cl = `^"` + in_cl + `"`
	default:
		cl = ""
		gol = "Object"
	}
	return
}
