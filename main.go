package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/twolodzko/goal/repl"
)

const (
	inputPrompt  string = ">"
	outputPrompt string = "=>"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Press ^C to exit.")
	fmt.Println()

	for {
		fmt.Printf("%s ", inputPrompt)

		input, err := repl.Read(reader)

		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		} else if len(input) > 0 {
			fmt.Printf("%s %s\n", outputPrompt, input)
		}
	}
}
