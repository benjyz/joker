package imports

import (
	"fmt"
	. "go/ast"
	"sort"
	. "strings"
)

/* Represents an 'import ( foo "bar/bletch/foo" )' line to be produced. */
type Import struct {
	Local       string // "foo", "_", ".", or empty
	LocalRef    string // local unless empty, in which case final component of full (e.g. "foo")
	Full        string // "bar/bletch/foo"
	substituted bool   // Had to substitute a different local name
}

/* Maps relative package (unix-style) names to their imports, non-emptiness, etc. */
type Imports struct {
	LocalNames map[string]string  // "foo" -> "bar/bletch/foo"; no "_" nor "." entries here
	FullNames  map[string]*Import // "bar/bletch/foo" -> ["foo", "bar/bletch/foo"]
}

/* Given desired local and the full (though relative) name of the
/* package, make sure the local name agrees with any existing entry
/* and isn't already used (picking an alternate local name if
/* necessary), add the mapping if necessary, and return the (possibly
/* alternate) local name. */
func AddImport(imports *Imports, local, full string, okToSubstitute bool) string {
	components := Split(full, "/")
	if e, found := imports.FullNames[full]; found {
		if e.Local == local {
			return e.LocalRef
		}
		if okToSubstitute {
			return e.LocalRef
		}
		panic(fmt.Sprintf("addImport(%s,%s) told to to replace (%s,%s)", local, full, e.Local, e.Full))
	}

	substituted := false
	localRef := local
	if local == "" {
		localRef = components[len(components)-1]
	}
	if localRef != "." {
		prevComponentIndex := len(components) - 1
		for {
			origLocalRef := localRef
			curFull, found := imports.LocalNames[localRef]
			if !found {
				break
			}
			substituted = true
			prevComponentIndex--
			if prevComponentIndex >= 0 {
				localRef = components[prevComponentIndex] + "_" + localRef
				continue
			} else if prevComponentIndex > -99 /* avoid infinite loop */ {
				localRef = fmt.Sprintf("%s_%d", origLocalRef, -prevComponentIndex)
				continue
			}
			panic(fmt.Sprintf("addImport(%s,%s) trying to replace (%s,%s)", localRef, full, imports.FullNames[curFull].LocalRef, curFull))
		}
		if imports.LocalNames == nil {
			imports.LocalNames = map[string]string{}
		}
		imports.LocalNames[localRef] = full
	}
	if imports.FullNames == nil {
		imports.FullNames = map[string]*Import{}
	}
	imports.FullNames[full] = &Import{local, localRef, full, substituted}
	return localRef
}

func SortedOriginalPackageImports(p *Package, f func(k string)) {
	imports := map[string]struct{}{}
	for _, f := range p.Files {
		for _, impSpec := range f.Imports {
			imports[impSpec.Path.Value] = struct{}{}
		}
	}
	var sortedImports []string
	for k, _ := range imports {
		sortedImports = append(sortedImports, k)
	}
	sort.Strings(sortedImports)
	for _, imp := range sortedImports {
		f(imp)
	}
}

func sortedImports(pi *Imports, f func(k string, v *Import)) {
	var keys []string
	for k, _ := range pi.FullNames {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := pi.FullNames[k]
		f(k, v)
	}
}

func QuotedImportList(pi *Imports, prefix string) string {
	imports := ""
	sortedImports(pi,
		func(k string, v *Import) {
			if v.Local == "" && !v.substituted {
				imports += prefix + `"` + k + `"`
			} else {
				imports += prefix + v.LocalRef + ` "` + k + `"`
			}
		})
	return imports
}
