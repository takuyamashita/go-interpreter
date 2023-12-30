package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/takuyamashita/go-interpreter/lexer"
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

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, progrom.String())
		io.WriteString(out, "\n")

	}
}

func printParserErrors(out io.Writer, errors []string) {

	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
