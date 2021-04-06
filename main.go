package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/twolodzko/gol/evaluator"
	"github.com/twolodzko/gol/parser"
	"github.com/twolodzko/gol/repl"
)

const prompt string = "> "

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			printHelp()
			return
		}
		evalScript()
		return
	}

	startRepl()
}

func printHelp() {
	fmt.Printf("%s [script]\n", os.Args[0])
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s             start REPL\n", os.Args[0])
	fmt.Printf("  %s script.lsp  evaluate script.lsp\n", os.Args[0])
	fmt.Printf("  %s -h,--help   display help\n", os.Args[0])
}

func evalScript() {
	code, err := parser.ReadFile(os.Args[1])
	if err != nil {
		log.Panic(err)
	}

	e := evaluator.NewEvaluator()
	objs, err := e.EvalString(code)
	if err != nil {
		log.Panic(err)
	}
	if len(objs) > 0 {
		fmt.Printf("%v\n", objs[len(objs)-1])
	}
}

func startRepl() {
	repl := repl.NewRepl(os.Stdin)

	fmt.Println("Press ^C to exit.")
	fmt.Println()

	for {
		fmt.Printf("%s", prompt)

		objs, err := repl.Repl()

		if err != nil {
			print(fmt.Sprintf("ERROR: %s", err))
			continue
		}

		for _, obj := range objs {
			print(fmt.Sprintf("%v", obj))
		}
	}
}

func print(msg string) {
	io.WriteString(os.Stdout, fmt.Sprintf("%s\n", msg))
}
