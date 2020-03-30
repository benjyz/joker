// This file is generated by generate-std.joke script. Do not edit manually!

// +build !gen_code

package string

import (
	"fmt"
	. "github.com/candid82/joker/core"
	"os"
)

func InternsOrThunks() {
	if VerbosityLevel > 0 {
		fmt.Fprintln(os.Stderr, "Lazily running fast version of string.InternsOrThunks().")
	}
	STD_thunk_string_isblank__var = __isblank_
	STD_thunk_string_capitalize__var = __capitalize_
	STD_thunk_string_isends_with__var = __isends_with_
	STD_thunk_string_escape__var = __escape_
	STD_thunk_string_isincludes__var = __isincludes_
	STD_thunk_string_index_of__var = __index_of_
	STD_thunk_string_join__var = __join_
	STD_thunk_string_last_index_of__var = __last_index_of_
	STD_thunk_string_lower_case__var = __lower_case_
	STD_thunk_string_pad_left__var = __pad_left_
	STD_thunk_string_pad_right__var = __pad_right_
	STD_thunk_string_re_quote__var = __re_quote_
	STD_thunk_string_replace__var = __replace_
	STD_thunk_string_replace_first__var = __replace_first_
	STD_thunk_string_reverse__var = __reverse_
	STD_thunk_string_split__var = __split_
	STD_thunk_string_split_lines__var = __split_lines_
	STD_thunk_string_isstarts_with__var = __isstarts_with_
	STD_thunk_string_trim__var = __trim_
	STD_thunk_string_trim_left__var = __trim_left_
	STD_thunk_string_trim_newline__var = __trim_newline_
	STD_thunk_string_trim_right__var = __trim_right_
	STD_thunk_string_triml__var = __triml_
	STD_thunk_string_trimr__var = __trimr_
	STD_thunk_string_upper_case__var = __upper_case_
}
