package jtypes

// Info on Joker types, including map of Joker type names to said type
// info.  A Joker type name is either unqualified (built-in, not
// namespace-rooted) or fully qualified by a namespace name
// (e.g. "go.std.example/SomeType").

type Info struct {
	ArgExtractFunc     string
	ArgClojureArgType  string // Clojure argument type for a Go function arg with my type
	ConvertFromClojure string // Pattern to convert a (scalar) %s to this type
	ConvertToClojure   string // Pattern to convert this type to an appropriate Clojure object
	AsJokerObject      string // Pattern to convert this type to a normal Joker type; empty string means wrap in a GoObject
}

var Nil = Info{}

var Error = Info{
	ArgExtractFunc:    "Error",
	ArgClojureArgType: "Error",
	ConvertToClojure:  "Error(%s%s)",
	AsJokerObject:     "Error(%s%s)",
}

var Bool = Info{
	ArgExtractFunc:    "Boolean",
	ArgClojureArgType: "Boolean",
	ConvertToClojure:  "Boolean(%s%s)",
	AsJokerObject:     "Boolean(%s%s)",
}

var Byte = Info{
	ArgExtractFunc:    "Byte",
	ArgClojureArgType: "Int",
	ConvertToClojure:  "Int(int(%s)%s)",
	AsJokerObject:     "Int(int(%s)%s)",
}

var Rune = Info{
	ArgExtractFunc:    "Char",
	ArgClojureArgType: "Char",
	ConvertToClojure:  "Char(%s%s)",
	AsJokerObject:     "Char(%s%s)",
}

var String = Info{
	ArgExtractFunc:    "String",
	ArgClojureArgType: "String",
	ConvertToClojure:  "String(%s%s)",
	AsJokerObject:     "String(%s%s)",
}

var Int = Info{
	ArgExtractFunc:    "Int",
	ArgClojureArgType: "Int",
	ConvertToClojure:  "Int(%s%s)",
	AsJokerObject:     "Int(%s%s)",
}

var Int32 = Info{
	ArgExtractFunc:    "Int32",
	ArgClojureArgType: "Int",
	ConvertToClojure:  "Int(int(%s)%s)",
	AsJokerObject:     "Int(int(%s)%s)",
}

var Int64 = Info{
	ArgExtractFunc:    "Int64",
	ArgClojureArgType: "Number",
	ConvertToClojure:  "BigInt(%s%s)",
	AsJokerObject:     "BigInt(%s%s)",
}

var UInt = Info{
	ArgExtractFunc:    "Uint",
	ArgClojureArgType: "Number",
	ConvertToClojure:  "BigIntU(uint64(%s)%s)",
	AsJokerObject:     "BigIntU(uint64(%s)%s)",
}

var UInt8 = Info{
	ArgExtractFunc:    "Uint8",
	ArgClojureArgType: "Int",
	ConvertToClojure:  "Int(int(%s)%s)",
	AsJokerObject:     "Int(int(%s)%s)",
}

var UInt16 = Info{
	ArgExtractFunc:    "Uint16",
	ArgClojureArgType: "Int",
	ConvertToClojure:  "Int(int(%s)%s)",
	AsJokerObject:     "Int(int(%s)%s)",
}

var UInt32 = Info{
	ArgExtractFunc:    "Uint32",
	ArgClojureArgType: "Number",
	ConvertToClojure:  "BigIntU(uint64(%s)%s)",
	AsJokerObject:     "BigIntU(uint64(%s)%s)",
}

var UInt64 = Info{
	ArgExtractFunc:    "Uint64",
	ArgClojureArgType: "Number",
	ConvertToClojure:  "BigIntU(%s%s)",
	AsJokerObject:     "BigIntU(%s%s)",
}

var UIntPtr = Info{
	ArgExtractFunc:    "UintPtr",
	ArgClojureArgType: "Number",
	AsJokerObject:     "Number(%s%s)",
}

var Float32 = Info{
	ArgExtractFunc:    "ABEND007(find these)",
	ArgClojureArgType: "Double",
	AsJokerObject:     "Double(float64(%s)%s)",
}

var Float64 = Info{
	ArgExtractFunc:    "ABEND007(find these)",
	ArgClojureArgType: "Double",
	AsJokerObject:     "Double(%s%s)",
}

var Complex128 = Info{
	ArgExtractFunc:    "ABEND007(find these)",
	ArgClojureArgType: "ABEND007(find these)",
}

var TypeMap = map[string]*Info{
	"Nil":        &Nil,
	"Error":      &Error,
	"Bool":       &Bool,
	"Byte":       &Byte,
	"Rune":       &Rune,
	"String":     &String,
	"Int":        &Int,
	"Int32":      &Int32,
	"Int64":      &Int64,
	"UInt":       &UInt,
	"UInt8":      &UInt8,
	"UInt16":     &UInt16,
	"UInt32":     &UInt32,
	"UInt64":     &UInt64,
	"UIntPtr":    &UIntPtr,
	"Float32":    &Float32,
	"Float64":    &Float64,
	"Complex128": &Complex128,
}
