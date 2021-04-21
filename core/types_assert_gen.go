// Generated by gen_types. Don't modify manually!

package core

import (
	"io"
)

func EnsureObjectIsComparable(obj Object, pattern string) Comparable {
	if c, yes := obj.(Comparable); yes {
		return c
	}
	panic(FailObject(obj, "Comparable", pattern))
}

func EnsureArgIsComparable(args []Object, index int) Comparable {
	obj := args[index]
	if c, yes := obj.(Comparable); yes {
		return c
	}
	panic(FailArg(obj, "Comparable", index))
}

func EnsureObjectIsVector(obj Object, pattern string) *Vector {
	if c, yes := obj.(*Vector); yes {
		return c
	}
	panic(FailObject(obj, "Vector", pattern))
}

func EnsureArgIsVector(args []Object, index int) *Vector {
	obj := args[index]
	if c, yes := obj.(*Vector); yes {
		return c
	}
	panic(FailArg(obj, "Vector", index))
}

func EnsureObjectIsChar(obj Object, pattern string) Char {
	if c, yes := obj.(Char); yes {
		return c
	}
	panic(FailObject(obj, "Char", pattern))
}

func EnsureArgIsChar(args []Object, index int) Char {
	obj := args[index]
	if c, yes := obj.(Char); yes {
		return c
	}
	panic(FailArg(obj, "Char", index))
}

func EnsureObjectIsString(obj Object, pattern string) String {
	if c, yes := obj.(String); yes {
		return c
	}
	panic(FailObject(obj, "String", pattern))
}

func EnsureArgIsString(args []Object, index int) String {
	obj := args[index]
	if c, yes := obj.(String); yes {
		return c
	}
	panic(FailArg(obj, "String", index))
}

func EnsureObjectIsSymbol(obj Object, pattern string) Symbol {
	if c, yes := obj.(Symbol); yes {
		return c
	}
	panic(FailObject(obj, "Symbol", pattern))
}

func EnsureArgIsSymbol(args []Object, index int) Symbol {
	obj := args[index]
	if c, yes := obj.(Symbol); yes {
		return c
	}
	panic(FailArg(obj, "Symbol", index))
}

func EnsureObjectIsKeyword(obj Object, pattern string) Keyword {
	if c, yes := obj.(Keyword); yes {
		return c
	}
	panic(FailObject(obj, "Keyword", pattern))
}

func EnsureArgIsKeyword(args []Object, index int) Keyword {
	obj := args[index]
	if c, yes := obj.(Keyword); yes {
		return c
	}
	panic(FailArg(obj, "Keyword", index))
}

func EnsureObjectIsRegex(obj Object, pattern string) *Regex {
	if c, yes := obj.(*Regex); yes {
		return c
	}
	panic(FailObject(obj, "Regex", pattern))
}

func EnsureArgIsRegex(args []Object, index int) *Regex {
	obj := args[index]
	if c, yes := obj.(*Regex); yes {
		return c
	}
	panic(FailArg(obj, "Regex", index))
}

func EnsureObjectIsBoolean(obj Object, pattern string) Boolean {
	if c, yes := obj.(Boolean); yes {
		return c
	}
	panic(FailObject(obj, "Boolean", pattern))
}

func EnsureArgIsBoolean(args []Object, index int) Boolean {
	obj := args[index]
	if c, yes := obj.(Boolean); yes {
		return c
	}
	panic(FailArg(obj, "Boolean", index))
}

func EnsureObjectIsTime(obj Object, pattern string) Time {
	if c, yes := obj.(Time); yes {
		return c
	}
	panic(FailObject(obj, "Time", pattern))
}

func EnsureArgIsTime(args []Object, index int) Time {
	obj := args[index]
	if c, yes := obj.(Time); yes {
		return c
	}
	panic(FailArg(obj, "Time", index))
}

func EnsureObjectIsNumber(obj Object, pattern string) Number {
	if c, yes := obj.(Number); yes {
		return c
	}
	panic(FailObject(obj, "Number", pattern))
}

func EnsureArgIsNumber(args []Object, index int) Number {
	obj := args[index]
	if c, yes := obj.(Number); yes {
		return c
	}
	panic(FailArg(obj, "Number", index))
}

func EnsureObjectIsSeqable(obj Object, pattern string) Seqable {
	if c, yes := obj.(Seqable); yes {
		return c
	}
	panic(FailObject(obj, "Seqable", pattern))
}

func EnsureArgIsSeqable(args []Object, index int) Seqable {
	obj := args[index]
	if c, yes := obj.(Seqable); yes {
		return c
	}
	panic(FailArg(obj, "Seqable", index))
}

func EnsureObjectIsCallable(obj Object, pattern string) Callable {
	if c, yes := obj.(Callable); yes {
		return c
	}
	panic(FailObject(obj, "Callable", pattern))
}

func EnsureArgIsCallable(args []Object, index int) Callable {
	obj := args[index]
	if c, yes := obj.(Callable); yes {
		return c
	}
	panic(FailArg(obj, "Callable", index))
}

func EnsureObjectIsType(obj Object, pattern string) *Type {
	if c, yes := obj.(*Type); yes {
		return c
	}
	panic(FailObject(obj, "Type", pattern))
}

func EnsureArgIsType(args []Object, index int) *Type {
	obj := args[index]
	if c, yes := obj.(*Type); yes {
		return c
	}
	panic(FailArg(obj, "Type", index))
}

func EnsureObjectIsMeta(obj Object, pattern string) Meta {
	if c, yes := obj.(Meta); yes {
		return c
	}
	panic(FailObject(obj, "Meta", pattern))
}

func EnsureArgIsMeta(args []Object, index int) Meta {
	obj := args[index]
	if c, yes := obj.(Meta); yes {
		return c
	}
	panic(FailArg(obj, "Meta", index))
}

func EnsureObjectIsInt(obj Object, pattern string) Int {
	if c, yes := obj.(Int); yes {
		return c
	}
	panic(FailObject(obj, "Int", pattern))
}

func EnsureArgIsInt(args []Object, index int) Int {
	obj := args[index]
	if c, yes := obj.(Int); yes {
		return c
	}
	panic(FailArg(obj, "Int", index))
}

func EnsureObjectIsDouble(obj Object, pattern string) Double {
	if c, yes := obj.(Double); yes {
		return c
	}
	panic(FailObject(obj, "Double", pattern))
}

func EnsureArgIsDouble(args []Object, index int) Double {
	obj := args[index]
	if c, yes := obj.(Double); yes {
		return c
	}
	panic(FailArg(obj, "Double", index))
}

func EnsureObjectIsStack(obj Object, pattern string) Stack {
	if c, yes := obj.(Stack); yes {
		return c
	}
	panic(FailObject(obj, "Stack", pattern))
}

func EnsureArgIsStack(args []Object, index int) Stack {
	obj := args[index]
	if c, yes := obj.(Stack); yes {
		return c
	}
	panic(FailArg(obj, "Stack", index))
}

func EnsureObjectIsMap(obj Object, pattern string) Map {
	if c, yes := obj.(Map); yes {
		return c
	}
	panic(FailObject(obj, "Map", pattern))
}

func EnsureArgIsMap(args []Object, index int) Map {
	obj := args[index]
	if c, yes := obj.(Map); yes {
		return c
	}
	panic(FailArg(obj, "Map", index))
}

func EnsureObjectIsSet(obj Object, pattern string) Set {
	if c, yes := obj.(Set); yes {
		return c
	}
	panic(FailObject(obj, "Set", pattern))
}

func EnsureArgIsSet(args []Object, index int) Set {
	obj := args[index]
	if c, yes := obj.(Set); yes {
		return c
	}
	panic(FailArg(obj, "Set", index))
}

func EnsureObjectIsAssociative(obj Object, pattern string) Associative {
	if c, yes := obj.(Associative); yes {
		return c
	}
	panic(FailObject(obj, "Associative", pattern))
}

func EnsureArgIsAssociative(args []Object, index int) Associative {
	obj := args[index]
	if c, yes := obj.(Associative); yes {
		return c
	}
	panic(FailArg(obj, "Associative", index))
}

func EnsureObjectIsReversible(obj Object, pattern string) Reversible {
	if c, yes := obj.(Reversible); yes {
		return c
	}
	panic(FailObject(obj, "Reversible", pattern))
}

func EnsureArgIsReversible(args []Object, index int) Reversible {
	obj := args[index]
	if c, yes := obj.(Reversible); yes {
		return c
	}
	panic(FailArg(obj, "Reversible", index))
}

func EnsureObjectIsNamed(obj Object, pattern string) Named {
	if c, yes := obj.(Named); yes {
		return c
	}
	panic(FailObject(obj, "Named", pattern))
}

func EnsureArgIsNamed(args []Object, index int) Named {
	obj := args[index]
	if c, yes := obj.(Named); yes {
		return c
	}
	panic(FailArg(obj, "Named", index))
}

func EnsureObjectIsComparator(obj Object, pattern string) Comparator {
	if c, yes := obj.(Comparator); yes {
		return c
	}
	panic(FailObject(obj, "Comparator", pattern))
}

func EnsureArgIsComparator(args []Object, index int) Comparator {
	obj := args[index]
	if c, yes := obj.(Comparator); yes {
		return c
	}
	panic(FailArg(obj, "Comparator", index))
}

func EnsureObjectIsRatio(obj Object, pattern string) *Ratio {
	if c, yes := obj.(*Ratio); yes {
		return c
	}
	panic(FailObject(obj, "Ratio", pattern))
}

func EnsureArgIsRatio(args []Object, index int) *Ratio {
	obj := args[index]
	if c, yes := obj.(*Ratio); yes {
		return c
	}
	panic(FailArg(obj, "Ratio", index))
}

func EnsureObjectIsBigFloat(obj Object, pattern string) *BigFloat {
	if c, yes := obj.(*BigFloat); yes {
		return c
	}
	panic(FailObject(obj, "BigFloat", pattern))
}

func EnsureArgIsBigFloat(args []Object, index int) *BigFloat {
	obj := args[index]
	if c, yes := obj.(*BigFloat); yes {
		return c
	}
	panic(FailArg(obj, "BigFloat", index))
}

func EnsureObjectIsNamespace(obj Object, pattern string) *Namespace {
	if c, yes := obj.(*Namespace); yes {
		return c
	}
	panic(FailObject(obj, "Namespace", pattern))
}

func EnsureArgIsNamespace(args []Object, index int) *Namespace {
	obj := args[index]
	if c, yes := obj.(*Namespace); yes {
		return c
	}
	panic(FailArg(obj, "Namespace", index))
}

func EnsureObjectIsVar(obj Object, pattern string) *Var {
	if c, yes := obj.(*Var); yes {
		return c
	}
	panic(FailObject(obj, "Var", pattern))
}

func EnsureArgIsVar(args []Object, index int) *Var {
	obj := args[index]
	if c, yes := obj.(*Var); yes {
		return c
	}
	panic(FailArg(obj, "Var", index))
}

func EnsureObjectIsFn(obj Object, pattern string) *Fn {
	if c, yes := obj.(*Fn); yes {
		return c
	}
	panic(FailObject(obj, "Fn", pattern))
}

func EnsureArgIsFn(args []Object, index int) *Fn {
	obj := args[index]
	if c, yes := obj.(*Fn); yes {
		return c
	}
	panic(FailArg(obj, "Fn", index))
}

func EnsureObjectIsDeref(obj Object, pattern string) Deref {
	if c, yes := obj.(Deref); yes {
		return c
	}
	panic(FailObject(obj, "Deref", pattern))
}

func EnsureArgIsDeref(args []Object, index int) Deref {
	obj := args[index]
	if c, yes := obj.(Deref); yes {
		return c
	}
	panic(FailArg(obj, "Deref", index))
}

func EnsureObjectIsAtom(obj Object, pattern string) *Atom {
	if c, yes := obj.(*Atom); yes {
		return c
	}
	panic(FailObject(obj, "Atom", pattern))
}

func EnsureArgIsAtom(args []Object, index int) *Atom {
	obj := args[index]
	if c, yes := obj.(*Atom); yes {
		return c
	}
	panic(FailArg(obj, "Atom", index))
}

func EnsureObjectIsRef(obj Object, pattern string) Ref {
	if c, yes := obj.(Ref); yes {
		return c
	}
	panic(FailObject(obj, "Ref", pattern))
}

func EnsureArgIsRef(args []Object, index int) Ref {
	obj := args[index]
	if c, yes := obj.(Ref); yes {
		return c
	}
	panic(FailArg(obj, "Ref", index))
}

func EnsureObjectIsKVReduce(obj Object, pattern string) KVReduce {
	if c, yes := obj.(KVReduce); yes {
		return c
	}
	panic(FailObject(obj, "KVReduce", pattern))
}

func EnsureArgIsKVReduce(args []Object, index int) KVReduce {
	obj := args[index]
	if c, yes := obj.(KVReduce); yes {
		return c
	}
	panic(FailArg(obj, "KVReduce", index))
}

func EnsureObjectIsPending(obj Object, pattern string) Pending {
	if c, yes := obj.(Pending); yes {
		return c
	}
	panic(FailObject(obj, "Pending", pattern))
}

func EnsureArgIsPending(args []Object, index int) Pending {
	obj := args[index]
	if c, yes := obj.(Pending); yes {
		return c
	}
	panic(FailArg(obj, "Pending", index))
}

func EnsureObjectIsFile(obj Object, pattern string) *File {
	if c, yes := obj.(*File); yes {
		return c
	}
	panic(FailObject(obj, "File", pattern))
}

func EnsureArgIsFile(args []Object, index int) *File {
	obj := args[index]
	if c, yes := obj.(*File); yes {
		return c
	}
	panic(FailArg(obj, "File", index))
}

func EnsureObjectIsio_Reader(obj Object, pattern string) io.Reader {
	if c, yes := obj.(io.Reader); yes {
		return c
	}
	panic(FailObject(obj, "io.Reader", pattern))
}

func EnsureArgIsio_Reader(args []Object, index int) io.Reader {
	obj := args[index]
	if c, yes := obj.(io.Reader); yes {
		return c
	}
	panic(FailArg(obj, "io.Reader", index))
}

func EnsureObjectIsio_Writer(obj Object, pattern string) io.Writer {
	if c, yes := obj.(io.Writer); yes {
		return c
	}
	panic(FailObject(obj, "io.Writer", pattern))
}

func EnsureArgIsio_Writer(args []Object, index int) io.Writer {
	obj := args[index]
	if c, yes := obj.(io.Writer); yes {
		return c
	}
	panic(FailArg(obj, "io.Writer", index))
}

func EnsureObjectIsStringReader(obj Object, pattern string) StringReader {
	if c, yes := obj.(StringReader); yes {
		return c
	}
	panic(FailObject(obj, "StringReader", pattern))
}

func EnsureArgIsStringReader(args []Object, index int) StringReader {
	obj := args[index]
	if c, yes := obj.(StringReader); yes {
		return c
	}
	panic(FailArg(obj, "StringReader", index))
}

func EnsureObjectIsio_RuneReader(obj Object, pattern string) io.RuneReader {
	if c, yes := obj.(io.RuneReader); yes {
		return c
	}
	panic(FailObject(obj, "io.RuneReader", pattern))
}

func EnsureArgIsio_RuneReader(args []Object, index int) io.RuneReader {
	obj := args[index]
	if c, yes := obj.(io.RuneReader); yes {
		return c
	}
	panic(FailArg(obj, "io.RuneReader", index))
}

func EnsureObjectIsChannel(obj Object, pattern string) *Channel {
	if c, yes := obj.(*Channel); yes {
		return c
	}
	panic(FailObject(obj, "Channel", pattern))
}

func EnsureArgIsChannel(args []Object, index int) *Channel {
	obj := args[index]
	if c, yes := obj.(*Channel); yes {
		return c
	}
	panic(FailArg(obj, "Channel", index))
}

func EnsureObjectIsGoObject(obj Object, pattern string) GoObject {
	if c, yes := obj.(GoObject); yes {
		return c
	}
	panic(FailObject(obj, "GoObject", pattern))
}

func EnsureArgIsGoObject(args []Object, index int) GoObject {
	obj := args[index]
	if c, yes := obj.(GoObject); yes {
		return c
	}
	panic(FailArg(obj, "GoObject", index))
}
