package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/takuyamashita/go-interpreter/evaluator"
	"github.com/takuyamashita/go-interpreter/lexer"
	"github.com/takuyamashita/go-interpreter/object"
	"github.com/takuyamashita/go-interpreter/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {

	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		progrom := p.ParseProgram()
		env := object.NewEnvironment()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(progrom, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParserErrors(out io.Writer, errors []string) {

	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
