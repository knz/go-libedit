// +build darwin freebsd linux openbsd netbsd dragonfly

package libedit

import (
	"fmt"
	"io"
	"syscall"
	"unsafe"
)

// - NetBSD and OpenBSD bundle libedit and place headers in /usr/include/readline.
// - FreeBSD bundles libedit and places headers in /usr/include/edit/readline.
// - OSX bundles libedit and places headers in /usr/include/editline.
// - DragonflyBSD bundles libedit and places headers in /usr/include/priv/readline.
// - Also both DragonflyBSD and OSX uses non-standard typedefs.
// - For Linux we use the bundled sources.

// #cgo openbsd netbsd freebsd dragonfly darwin LDFLAGS: -ledit
// #cgo openbsd netbsd freebsd dragonfly darwin CPPFLAGS: -Ishim
// #cgo openbsd netbsd   CPPFLAGS: -I/usr/include/readline
// #cgo darwin           CPPFLAGS: -I/usr/include/editline
// #cgo freebsd          CPPFLAGS: -I/usr/include/edit/readline
// #cgo dragonfly        CPPFLAGS: -I/usr/include/priv/readline
// #cgo dragonfly darwin CPPFLAGS: -Dweird_completion_typedef
// #cgo linux          LDFLAGS: -lncurses
// #cgo linux          CFLAGS: -Wno-unused-result
// #cgo linux          CPPFLAGS: -Isrc -Isrc/c-libedit -Isrc/c-libedit/editline -Isrc/c-libedit/linux-build
//
// #include <readline.h>
// #include <stdio.h>
// #include <stdlib.h>
//
// // Some helper functions that make the Go code easier on the eye.
// void go_libedit_printstring(char *s) { printf("%s\n", s); }
// void go_libedit_set_string_array(char **ar, int p, char *s) { ar[p] = s; }
//
// // This function is defined in edit_unix_completion.go.
// extern char **go_libedit_autocomplete(char *word, char *line, int start, int end);
//
// // This function adds the const qualifier which the Go //export directive can't generate.
// static char **wrap_autocomplete(const char *word, int start, int end) {
//     return go_libedit_autocomplete((char*)word, rl_line_buffer, start, end);
// }
//
// #ifdef weird_completion_typedef
// #define go_libedit_completion_func_t CPPFunction
// #else
// #define go_libedit_completion_func_t rl_completion_func_t
// #endif
// go_libedit_completion_func_t *go_libedit_autocomplete_ptr = wrap_autocomplete;
import "C"

var cAppName *C.char

func Initialize(appname string) error {
	if appname != "" {
		// rl_readline_name allows a user to customize their
		// ~/.editrc configuration per-app.
		cAppName = C.CString(appname)
		C.rl_readline_name = cAppName
	}

	r := C.rl_initialize()
	if r < 0 {
		Cleanup()
		return fmt.Errorf("unable to initialize libedit: %v", syscall.Errno(-r))
	}

	// rl_attempted_completion_function is called pre-completion to generate
	// a set of candidate strings.
	C.rl_attempted_completion_function = C.go_libedit_autocomplete_ptr
	// rl_attempted_completion_over, when non-zero, disables filename completion
	// after the attempted completion function has run.
	C.rl_attempted_completion_over = 1

	return nil
}

var histFile *C.char
var autoSaveHistory = false
var histExpand = true

func UseHistory(file string, autoSave bool, expand bool) error {
	newHistFile := C.CString(file)
	_, err := C.read_history(newHistFile)
	if err != nil && err != syscall.ENOENT {
		C.free(unsafe.Pointer(newHistFile))
		return err
	}

	if histFile != nil {
		C.free(unsafe.Pointer(histFile))
	}
	histFile = newHistFile
	autoSaveHistory = autoSave
	histExpand = expand
	return nil
}

func AddHistory(line string) error {
	cline := C.CString(line)
	defer C.free(unsafe.Pointer(cline))
	C.add_history(cline)
	if autoSaveHistory && histFile != nil {
		_, err := C.write_history(histFile)
		return err
	}
	return nil
}

var cPrompt *C.char

func SetPrompt(prompt string) {
	if cPrompt != nil {
		C.free(unsafe.Pointer(cPrompt))
	}
	cPrompt = C.CString(prompt)
}

func Readline() (string, error) {
	for {
		l, err := C.readline(cPrompt)
		if err != nil {
			if l != nil {
				C.free(unsafe.Pointer(l))
			}
			return "", err
		}
		if l == nil {
			return "", io.EOF
		}

		if histExpand {
			// Process history expansion commands.
			var exp *C.char
			res := C.history_expand(l, &exp)
			C.free(unsafe.Pointer(l))
			if res < 0 {
				// Input to history_expand was not valid; history_expand has
				// printed an error message already, we can just ask for a new
				// line.
				continue
			}
			if res == 2 {
				// History command says print result but do not execute.
				C.go_libedit_printstring(exp)
				C.free(unsafe.Pointer(exp))
				continue
			}
			l = exp
		}

		ret := C.GoString(l)
		C.free(unsafe.Pointer(l))
		return ret, nil
	}
}

func Cleanup() {
	if cPrompt != nil {
		C.free(unsafe.Pointer(cPrompt))
		cPrompt = nil
	}
	if histFile != nil {
		C.free(unsafe.Pointer(histFile))
		histFile = nil
	}
	if cAppName != nil {
		C.free(unsafe.Pointer(cAppName))
		cAppName = nil
	}
}
