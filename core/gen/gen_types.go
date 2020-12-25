package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type (
	TypeInfo struct {
		Name     string
		TypeName string
		ShowName string
	}
)

var header string = `// Generated by gen_types. Don't modify manually!

package core
`

var importFmt string = `
import (
	"fmt"
	"io"
)
`

var ensureObjectIsTemplate string = `
func EnsureObjectIs{{.Name}}(obj Object, pattern string) {{.TypeName}} {
	switch c := obj.(type) {
	case {{.TypeName}}:
		return c
	default:
		if pattern == "" {
			pattern = "%s"
		}
		msg := fmt.Sprintf("Expected %s, got %s", "{{.ShowName}}", obj.GetType().ToString(false))
		panic(RT.NewError(fmt.Sprintf(pattern, msg)))
	}
}
`

var ensureArgIsTemplate string = `
func EnsureArgIs{{.Name}}(args []Object, index int) {{.TypeName}} {
	switch c := args[index].(type) {
	case {{.TypeName}}:
		return c
	default:
		panic(RT.NewArgTypeError(index, c, "{{.ShowName}}"))
	}
}
`

var infoTemplate string = `
func (x {{.TypeName}}) WithInfo(info *ObjectInfo) Object {
	x.info = info
	return x
}
`

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func generateAssertions(types []string) {
	filename := "types_assert_gen.go"
	f, err := os.Create(filename)
	checkError(err)
	defer f.Close()

	var ensureObjectIs = template.Must(template.New("assert").Parse(ensureObjectIsTemplate))
	var ensureArgIs = template.Must(template.New("ensure").Parse(ensureArgIsTemplate))
	f.WriteString(header)
	f.WriteString(importFmt)
	for _, t := range types {
		typeInfo := TypeInfo{
			Name:     t,
			TypeName: t,
			ShowName: t,
		}
		if t[0] == '*' {
			typeInfo.Name = t[1:]
			typeInfo.ShowName = typeInfo.Name
		} else if strings.ContainsRune(t, '.') {
			typeInfo.Name = strings.ReplaceAll(t, ".", "_")
		}
		ensureObjectIs.Execute(f, typeInfo)
		ensureArgIs.Execute(f, typeInfo)
	}
}

func generateInfo(types []string) {
	filename := "types_info_gen.go"
	f, err := os.Create(filename)
	checkError(err)
	defer f.Close()

	var info = template.Must(template.New("info").Parse(infoTemplate))

	f.WriteString(header)
	for _, t := range types {
		typeInfo := TypeInfo{
			Name:     t,
			TypeName: t,
		}
		if t[0] == '*' {
			typeInfo.Name = t[1:]
		}
		info.Execute(f, typeInfo)
	}
}

func main() {
	cmd := os.Args[1]
	switch cmd {
	case "assert":
		generateAssertions(os.Args[2:])
	case "info":
		generateInfo(os.Args[2:])
	default:
		fmt.Println("Unknown command: ", cmd)
	}
}
