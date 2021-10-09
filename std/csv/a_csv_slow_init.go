// This file is generated by generate-std.joke script. Do not edit manually!


package csv

import (
	. "github.com/candid82/joker/core"
	"fmt"
	"os"
)

func InternsOrThunks() {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running slow version of csv.InternsOrThunks().")
	}
	csvNamespace.ResetMeta(MakeMeta(nil, `Reads and writes comma-separated values (CSV) files as defined in RFC 4180.`, "1.0"))

	
	csvNamespace.InternVar("csv-seq", csv_seq_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("rdr")), NewVectorFrom(MakeSymbol("rdr"), MakeSymbol("opts"))),
			`Returns the csv records from rdr as a lazy sequence.
  rdr must be a string or implement io.Reader.
  opts may have the following keys:

  :comma - field delimiter (defaults to ',').
  Must be a valid char and must not be \r, \n,
  or the Unicode replacement character (0xFFFD).

  :comment - comment character (defaults to 0 meaning no comments).
  Lines beginning with the comment character without preceding whitespace are ignored.
  With leading whitespace the comment character becomes part of the
  field, even if trim-leading-space is true.
  comment must be a valid chat and must not be \r, \n,
  or the Unicode replacement character (0xFFFD).
  It must also not be equal to comma.

  :fields-per-record - number of expected fields per record.
  If fields-per-record is positive, csv-seq requires each record to
  have the given number of fields. If fields-per-record is 0 (default), csv-seq sets it to
  the number of fields in the first record, so that future records must
  have the same field count. If fields-per-record is negative, no check is
  made and records may have a variable number of fields.

  :lazy-quotes - if true, a quote may appear in an unquoted field and a
  non-doubled quote may appear in a quoted field. Default value is false.

  :trim-leading-space - if true, leading white space in a field is ignored.
  This is done even if the field delimiter, comma, is white space.
  Default value is false.`, "1.0"))

	csvNamespace.InternVar("write", write_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("f"), MakeSymbol("data")), NewVectorFrom(MakeSymbol("f"), MakeSymbol("data"), MakeSymbol("opts"))),
			`Writes records to a CSV encoded file.
  f must be io.Writer (for example, as returned by joker.os/create).
  data must be Seqable, each element of which must be Seqable as well.
  opts is as in joker.csv/write-string.`, "1.0"))

	csvNamespace.InternVar("write-string", write_string_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("data")), NewVectorFrom(MakeSymbol("data"), MakeSymbol("opts"))),
			`Writes records to a string in CSV format and returns the string.
  data must be Seqable, each element of which must be Seqable as well.
  opts may have the following keys:

  :comma - field delimiter (defaults to ',')

  :use-crlf - if true, uses \r\n as the line terminator. Default value is false.`, "1.0").Plus(MakeKeyword("tag"), String{S: "String"}))

}
