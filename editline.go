package libedit

import "errors"

type EditLine int

type LeftPromptGenerator interface {
	GetLeftPrompt() string
}

type RightPromptGenerator interface {
	GetRightPrompt() string
}

type CompletionGenerator interface {
	GetCompletions(word string) []string
}

var ErrInterrupted = errors.New("interrupted")
