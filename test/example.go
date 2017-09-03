package main

import (
	"fmt"
	"io"
	"log"

	libedit "github.com/knz/go-libedit"
)

func main() {
	if err := libedit.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer libedit.Cleanup()

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
	}
	fmt.Println("goodbye!")
}
