// +build !darwin,!freebsd,!linux,!openbsd,!netbsd,!dragonfly

package libedit

import (
	"bufio"
	"os"
)

var reader *bufio.Reader
var prompt string

func Initialize() error {
	reader = bufio.NewReader(os.Stdin)
	return nil
}

func Cleanup() {}

func SetPrompt(p string) { prompt = p }

func Readline() (string, error) {
	os.Stdout.WriteString(prompt)
	os.Stdout.Sync()
	return reader.ReadString('\n')
}
