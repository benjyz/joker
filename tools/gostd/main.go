package main

import (
	"fmt"
	"github.com/candid82/joker/tools/gostd/abends"
	"github.com/candid82/joker/tools/gostd/godb"
	"github.com/candid82/joker/tools/gostd/paths"
	. "github.com/candid82/joker/tools/gostd/utils"
	"go/build"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const VERSION = "0.1"

/* Want to support e.g.:

     net/dial.go:DialTimeout(network, address string, timeout time.Duration) => _ Conn, _ error

   I.e. a function defined in one package refers to a type defined in
   another (a different directory, even).

   Sample routines include (from 'net' package):
     - lookupMX
     - queryEscape
   E.g.:
     ./gostd --dir $PWD/tests/small --output-code 2>&1 | grep -C20 lookupMX

*/

var goSourcePath string

func notOption(arg string) bool {
	return arg == "--" || arg == "-" || !strings.HasPrefix(arg, "-")
}

func usage() {
	fmt.Print(`
Usage: gostd options...

Options:
  --go <go-source-dir>        # Location of Golang's src/ subdirectory (defaults to build.Default.GOROOT)
  --others <package-spec>...  # Location of another package directory, or a file with one <package-spec> per line
  --go-source-path <path>     # Overrides $GOPATH/src/ as "root" of <package-spec> specifications
  --overwrite                 # Overwrite any existing <joker-std-subdir> files, leaving existing files intact
  --replace                   # 'rm -fr <joker-std-subdir>' before creating <joker-std-subdir>
  --fresh                     # (Default) Refuse to overwrite existing <joker-std-subdir> directory
  --joker <joker-source-dir>  # Modify pertinent source files to reflect packages being created
  --import-from <dir>         # Override <joker-source-dir> with <dir> for use in generated import decls; "--" means use default
  --verbose, -v               # Print info on what's going on
  --summary                   # Print summary of #s of types, functions, etc.
  --output-code               # Print generated code to stdout (used by test.sh)
  --empty                     # Generate empty packages (those with no Joker code)
  --dump                      # Use go's AST dump API on pertinent elements (functions, types, etc.)
  --no-timestamp              # Don't put the time (and version) info in generated/modified files
  --undo                      # Undo effects of --joker ...
  --help, -h                  # Print this information

If <joker-std-subdir> is not specified, no Go nor Clojure source files
(nor any other files nor directories) are created, effecting a sort of
"dry run".
`)
	os.Exit(0)
}

func listOfOthers(other string) (others []string) {
	o := filepath.Join(goSourcePath, other)
	s, e := os.Stat(o)
	if e != nil {
		o = other // try original without $GOPATH/src/ prefix
		s, e = os.Stat(o)
	}
	Check(e)
	if s.IsDir() {
		return []string{o}
	}
	fmt.Fprintf(os.Stderr, "files not yet supported: %s\n", other)
	os.Exit(3)
	return
}

func main() {
	godb.Fset = token.NewFileSet() // positions are relative to Fset

	length := len(os.Args)
	goSourceDir := ""
	goSourceDirVia := ""
	goSourcePath = os.Getenv("GOPATH")
	var others []string
	var otherSourceDirs []string
	jokerSourceDir := ""
	jokerImportDir := ""
	replace := false
	overwrite := false
	summary := false
	generateEmpty := false
	outputCode := false
	undo := false

	for i := 1; i < length; i++ { // shift
		a := os.Args[i]
		if a[0] == "-"[0] {
			switch a {
			case "--help", "-h":
				usage()
			case "--version", "-V":
				fmt.Printf("%s version %s\n", os.Args[0], VERSION)
			case "--no-timestamp":
				noTimeAndVersion = true
			case "--dump":
				godb.Dump = true
				WalkDump = true
			case "--overwrite":
				overwrite = true
				replace = false
			case "--replace":
				replace = true
				overwrite = false
			case "--fresh":
				replace = false
				overwrite = false
			case "--verbose", "-v":
				godb.Verbose = true
			case "--summary":
				summary = true
			case "--output-code":
				outputCode = true
			case "--empty":
				generateEmpty = true
			case "--undo":
				undo = true
			case "--go":
				if goSourceDir != "" {
					fmt.Fprintf(os.Stderr, "cannot specify --go <go-source-dir> more than once\n")
					os.Exit(1)
				}
				if i < length-1 && notOption(os.Args[i+1]) {
					i += 1 // shift
					goSourceDir = os.Args[i]
					goSourceDirVia = "--go"
				} else {
					fmt.Fprintf(os.Stderr, "missing path after --go option\n")
					os.Exit(1)
				}
			case "--others":
				if i >= length-1 || !notOption(os.Args[i+1]) {
					fmt.Fprintf(os.Stderr, "missing package-spec(s) after --others option\n")
					os.Exit(1)
				}
				for i < length-1 && notOption(os.Args[i+1]) {
					i += 1 // shift
					others = append(others, os.Args[i])
				}
			case "--go-source-path":
				if i < length-1 && notOption(os.Args[i+1]) {
					i += 1 // shift
					goSourcePath = os.Args[i]
				} else {
					fmt.Fprintf(os.Stderr, "missing package-spec(s) after --go-source-path option\n")
					os.Exit(1)
				}
			case "--joker":
				if jokerSourceDir != "" {
					fmt.Fprintf(os.Stderr, "cannot specify --joker <joker-source-dir> more than once\n")
					os.Exit(1)
				}
				if i < length-1 && notOption(os.Args[i+1]) {
					i += 1 // shift
					jokerSourceDir = os.Args[i]
				} else {
					fmt.Fprintf(os.Stderr, "missing path after --joker option\n")
					os.Exit(1)
				}
			case "--import-from":
				if jokerImportDir != "" {
					fmt.Fprintf(os.Stderr, "cannot specify --import-from <import-dir> more than once\n")
					os.Exit(1)
				}
				if i < length-1 && notOption(os.Args[i+1]) {
					i += 1 // shift
					jokerImportDir = os.Args[i]
				} else {
					fmt.Fprintf(os.Stderr, "missing path after --import-from option; got %s, which looks like an option\n", os.Args[i+1])
					os.Exit(1)
				}
			default:
				fmt.Fprintf(os.Stderr, "unrecognized option %s\n", a)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "extraneous argument(s) starting with: %s\n", a)
			os.Exit(1)
		}
	}

	if godb.Verbose {
		fmt.Printf("Default context: %v\n", build.Default)
	}

	if goSourceDir == "" {
		goSourceDir = build.Default.GOROOT
		goSourceDirVia = "build.Default.GOROOT"
	}

	goSrcDir := paths.NewNativePath(goSourceDir).Join("src")

	if fi, e := os.Stat(goSrcDir.Join("go").String()); e != nil || !fi.IsDir() {
		if m, e := filepath.Glob(goSrcDir.Join("*.go").String()); e != nil || m == nil || len(m) == 0 {
			fmt.Fprintf(os.Stderr, "Does not exist or is not a Go source directory: %s (specified via %s);\n%v",
				goSrcDir, goSourceDirVia, m)
			os.Exit(2)
		}
	}

	if goSourcePath == "" {
		fmt.Fprintf(os.Stderr, "no Go source path defined via either $GOPATH or --go-source-path")
		os.Exit(1)
	}
	if fi, e := os.Stat(goSourcePath); e == nil && fi.IsDir() && filepath.Base(goSourcePath) != "src" {
		goSourcePath = filepath.Join(goSourcePath, "src")
	}

	if jokerSourceDir != "" && jokerImportDir == "" {
		godb.SetJokerSourceDir(jokerSourceDir, goSourcePath)
	} else if jokerImportDir != "" && jokerImportDir != "--" {
		godb.SetJokerSourceDir(jokerImportDir, goSourcePath)
	}

	if godb.Verbose {
		fmt.Printf("goSrcDir: %s\n", goSrcDir)
		fmt.Printf("goSourcePath: %s\n", goSourcePath)
		fmt.Printf("JokerSourceDir (for import line): %s\n", godb.JokerSourceDir)
		for _, o := range others {
			otherSourceDirs = append(otherSourceDirs, listOfOthers(o)...)
		}
		for _, o := range otherSourceDirs {
			fmt.Printf("Other: %s\n", o)
		}
	}

	jokerLibDir := ""
	if jokerSourceDir != "" && jokerSourceDir != "-" {
		jokerLibDir = filepath.Join(jokerSourceDir, "std", "go", "std")
		if replace || undo {
			if e := os.RemoveAll(jokerLibDir); e != nil {
				fmt.Fprintf(os.Stderr, "Unable to effectively 'rm -fr %s'\n", jokerLibDir)
				os.Exit(1)
			}
		}

		if undo {
			RegisterPackages([]string{}, jokerSourceDir)
			RegisterJokerFiles([]string{}, jokerSourceDir)
			RegisterGoTypeSwitch([]TypeInfo{}, jokerSourceDir, false)
			os.Exit(0)
		}

		if !overwrite {
			if _, e := os.Stat(jokerLibDir); e == nil ||
				(!strings.Contains(e.Error(), "no such file or directory") &&
					!strings.Contains(e.Error(), "The system cannot find the")) { // Sometimes "...file specified"; other times "...path specified"!
				msg := "already exists"
				if e != nil {
					msg = e.Error()
				}
				fmt.Fprintf(os.Stderr, "Refusing to populate existing directory %s; please 'rm -fr' first, or specify --overwrite or --replace: %s\n",
					jokerLibDir, msg)
				os.Exit(1)
			}
			if e := os.MkdirAll(jokerLibDir, 0777); e != nil {
				fmt.Fprintf(os.Stderr, "Cannot 'mkdir -p %s': %s\n", jokerLibDir, e.Error())
				os.Exit(1)
			}
		}
	}

	godb.AddMapping(goSrcDir, "go.std.")
	root := goSrcDir.Join(".")
	AddWalkDir(goSrcDir, root, "go.std.")

	for _, o := range otherSourceDirs {
		op := paths.NewNativePath(o)
		root := op.Join(".")
		AddWalkDir(op, root, "x.y.z.")
	}

	err, badDir := WalkAllDirs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory %s: %v", badDir, err)
		os.Exit(1)
	}

	SortAllTypes()

	GenTypeInfo()

	SortedConstantInfoMap(GoConstants,
		func(c string, ci *ConstantInfo) {
			GenConstant(ci)
		})

	SortedVariableInfoMap(GoVariables,
		func(c string, ci *VariableInfo) {
			GenVariable(ci)
		})

	/* Generate type-code snippets in sorted order. */
	SortedTypeInfoMap(TypesByGoName(),
		func(t string, ti TypeInfo) {
			if ti == nil {
				panic(fmt.Sprintf("nil ti for `%s'", t))
			}
			if ti.TypeSpec() != nil {
				GenType(t, ti)
			}
		})

	/* Generate function-code snippets in alphabetical order. */
	SortedFuncInfoMap(QualifiedFunctions,
		func(f string, v *FuncInfo) {
			//			fmt.Printf("main.go: Qualifiedfunctions[%s]\n", f)
			if v.Fd != nil && v.Fd.Recv == nil {
				GenStandalone(v)
			} else {
				GenReceiver(v)
			}
		})

	OutputPackageCode(jokerLibDir, outputCode, generateEmpty)

	if jokerSourceDir != "" && jokerSourceDir != "-" {
		var packagesArray = []string{} // Relative package pathnames in alphabetical order
		var dotJokeArray = []string{}  // Relative package pathnames in alphabetical order

		SortedPackagesInfo(PackagesInfo,
			func(p string, i *PackageInfo) {
				if !generateEmpty && !i.NonEmpty {
					return
				}
				if i.HasGoFiles {
					packagesArray = append(packagesArray, p)
				}
				dotJokeArray = append(dotJokeArray, p)
			})
		RegisterPackages(packagesArray, jokerSourceDir)
		RegisterJokerFiles(dotJokeArray, jokerSourceDir)
	}

	RegisterGoTypeSwitch(AllTypesSorted(), jokerSourceDir, outputCode)

	if godb.Verbose || summary {
		fmt.Printf("ABENDs:")
		abends.PrintAbends()
		fmt.Printf(`
Totals: functions=%d generated=%d (%s%%)
          non-receivers=%d (%s%%) generated=%d (%s%%)
          receivers=%d (%s%%) generated=%d (%s%%)
          methods=%d (%s%%) generated=%d (%s%%)
        types=%d
          constructable=%d ctors=%d (%s%%)
        constants=%d generated=%d (%s%%)
        variables=%d generated=%d (%s%%)
`,
			NumFunctions, NumGeneratedFunctions, pct(NumGeneratedFunctions, NumFunctions),
			NumStandalones, pct(NumStandalones, NumFunctions), NumGeneratedStandalones, pct(NumGeneratedStandalones, NumStandalones),
			NumReceivers, pct(NumReceivers, NumFunctions), NumGeneratedReceivers, pct(NumGeneratedReceivers, NumReceivers),
			godb.NumMethods, pct(godb.NumMethods, NumFunctions), godb.NumGeneratedMethods, pct(godb.NumGeneratedMethods, godb.NumMethods),
			NumTypes,
			NumCtableTypes, NumGeneratedCtors, pct(NumGeneratedCtors, NumCtableTypes),
			NumConstants, NumGeneratedConstants, pct(NumGeneratedConstants, NumConstants),
			NumVariables, NumGeneratedVariables, pct(NumGeneratedVariables, NumVariables))
	}

	os.Exit(0)
}

func pct(i, j int) string {
	if j == 0 {
		return "--"
	}
	return fmt.Sprintf("%0.2f", (float64(i)/float64(j))*100.0)
}

func init() {
	nonEmptyLineRegexp = regexp.MustCompile(`(?m)^(.)`)
}
