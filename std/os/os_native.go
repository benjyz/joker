package os

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	. "github.com/candid82/joker/core"
)

func env() Object {
	res := EmptyArrayMap()
	for _, v := range os.Environ() {
		parts := strings.Split(v, "=")
		res.Add(String{S: parts[0]}, String{S: parts[1]})
	}
	return res
}

func commandArgs() Object {
	res := EmptyVector
	for _, arg := range os.Args {
		res = res.Conjoin(String{S: arg})
	}
	return res
}

func sh(name string, args []string) Object {
	cmd := exec.Command(name, args...)
	stdoutReader, err := cmd.StdoutPipe()
	PanicOnErr(err)
	stderrReader, err := cmd.StderrPipe()
	PanicOnErr(err)
	if err = cmd.Start(); err != nil {
		panic(RT.NewError(err.Error()))
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdoutReader)
	stdoutString := buf.String()
	buf = new(bytes.Buffer)
	buf.ReadFrom(stderrReader)
	stderrString := buf.String()
	err = cmd.Wait()
	res := EmptyArrayMap()
	res.Add(MakeKeyword("success"), Bool{B: err == nil})
	if err != nil {
		res.Add(MakeKeyword("err-msg"), String{S: err.Error()})
	}
	res.Add(MakeKeyword("out"), String{S: stdoutString})
	res.Add(MakeKeyword("err"), String{S: stderrString})
	return res
}

func mkdir(name string, perm int) Object {
	err := os.Mkdir(name, os.FileMode(perm))
	PanicOnErr(err)
	return NIL
}

func readDir(dirname string) Object {
	files, err := ioutil.ReadDir(dirname)
	PanicOnErr(err)
	res := EmptyVector
	name := MakeKeyword("name")
	size := MakeKeyword("size")
	mode := MakeKeyword("mode")
	isDir := MakeKeyword("dir?")
	modTime := MakeKeyword("modtime")
	for _, f := range files {
		m := EmptyArrayMap()
		m.Add(name, MakeString(f.Name()))
		m.Add(size, MakeInt(int(f.Size())))
		m.Add(mode, MakeInt(int(f.Mode())))
		m.Add(isDir, MakeBool(f.IsDir()))
		m.Add(modTime, MakeInt(int(f.ModTime().Unix())))
		res = res.Conjoin(m)
	}
	return res
}
