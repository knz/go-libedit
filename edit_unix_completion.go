// +build darwin freebsd linux openbsd netbsd dragonfly

package libedit

// #include <stdlib.h>
// #include "c_editline.h"
// typedef char* pchar;
import "C"

type AutoCompleter interface {
	Do(word, line string, start, end int) []string
}

var Completer AutoCompleter

//export go_libedit_autocomplete
func go_libedit_autocomplete(cWord, cLine *C.char, cStart, cEnd C.int) **C.char {
	// word is the current word being completed.
	// line is the entire line.
	// start, end are the position of the word in line being completed.

	if Completer == nil {
		return nil
	}

	line := C.GoString(cLine)
	word := line[int(cStart):int(cEnd)]
	matches := Completer.Do(word, line, int(cStart), int(cEnd))
	if len(matches) == 0 {
		return nil
	}

	array := (**C.char)(C.malloc(C.size_t(C.sizeof_pchar * (len(matches) + 1))))
	for i, m := range matches {
		C.go_libedit_set_string_array(array, C.int(i), C.CString(m))
	}
	C.go_libedit_set_string_array(array, C.int(len(matches)), nil)
	return array
}
