# go-libedit
Go wrapper around the BSD libedit replacement to GNU readline.

How to use:

- see `test/example.go` for a demo.
- basic idea: call `Initialize()` once. Then call `SetPrompt()` and
  `Readline()` as needed. Finally call `Cleanup()`.

