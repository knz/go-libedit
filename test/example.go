package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	libedit "github.com/knz/go-libedit"
)

type completeHello struct{}

func (_ completeHello) Do(word, line string, start, end int) []string {
	if strings.HasPrefix(word, "he") {
		return []string{"hello!"}
	}
	return nil
}

func main() {
	if err := libedit.Initialize("example"); err != nil {
		log.Fatalf("init: %v", err)
	}
	defer libedit.Cleanup()

	if err := libedit.UseHistory("", true, true); err != nil {
		log.Fatalf("use hist: %v", err)
	}

	libedit.Completer = completeHello{}

	libedit.SetPrompt("hello> ")
	for {
		s, err := libedit.Readline()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println("echo", s)
		if err := libedit.AddHistory(s); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("goodbye!")
}
