package core

import (
	"os"
	"unsafe"
)

type (
	File struct {
		*os.File
	}
)

func (f *File) ToString(escape bool) string {
	return "#object[File]"
}

func (f *File) TypeToString(escape bool) string {
	return f.GetType().ToString(escape)
}

func (f *File) Equals(other interface{}) bool {
	return f == other
}

func (f *File) GetInfo() *ObjectInfo {
	return nil
}

func (f *File) GetType() *Type {
	return TYPE.File
}

func (f *File) Hash() uint32 {
	return HashPtr(uintptr(unsafe.Pointer(f)))
}

func (f *File) WithInfo(info *ObjectInfo) Object {
	return f
}

// To satisfy Named interface
func (f *File) Namespace() string {
	return ""
}

func MakeFile(f *os.File) *File {
	return &File{f}
}

func ExtractFile(args []Object, index int) *File {
	return EnsureArgIsFile(args, index)
}
