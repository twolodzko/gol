package main

import (
	"fmt"
	"os"

	"github.com/twolodzko/goal/repl"
)

const (
	inputPrompt  string = ""
	outputPrompt string = "; "
)

func print(msg string) {
	fmt.Printf("%s%s\n", outputPrompt, msg)
}

func main() {
	fmt.Println("Press ^C to exit.")
	fmt.Println()

	for {
		fmt.Printf("%s", inputPrompt)

		input, err := repl.Read(os.Stdin)

		if err != nil {
			print(fmt.Sprintf("ERROR: %s", err))
		} else if len(input) > 0 {
			print(input)
		}
	}
}
