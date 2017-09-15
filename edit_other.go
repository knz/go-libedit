// +build !darwin,!freebsd,!linux,!openbsd,!netbsd,!dragonfly

package libedit

import (
	"bufio"
	"os"
)

var reader *bufio.Reader
var prompt string

func Initialize(_ string) error {
	reader = bufio.NewReader(os.Stdin)
	return nil
}

func UseHistory(_ string, _, _ bool) error { return nil }
func AddHistory(_ string) error            { return nil }
func Cleanup()                             {}
func SetPrompt(p string)                   { prompt = p }

func Readline() (string, error) {
	os.Stdout.WriteString(prompt)
	os.Stdout.Sync()
	return reader.ReadString('\n')
}
