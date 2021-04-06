package main

import (
	"fmt"
	"io"
	"os"

	"github.com/twolodzko/gol/repl"
)

const (
	inputPrompt  string = "> "
	outputPrompt string = "=> "
)

func print(msg string) {
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s\n", outputPrompt, msg))
}

func main() {
	repl := repl.NewRepl(os.Stdin)

	fmt.Println("Press ^C to exit.")
	fmt.Println()

	for {
		fmt.Printf("%s", inputPrompt)

		out, err := repl.Repl()

		if err != nil {
			print(fmt.Sprintf("ERROR: %s", err))
			continue
		}

		for _, obj := range out {
			print(fmt.Sprintf("%v", obj))
		}
	}
}
