package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/takuyamashita/go-interpreter/lexer"
	"github.com/takuyamashita/go-interpreter/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {

	scanner := bufio.NewScanner(in)

	for {
		// Print the prompt.
		fmt.Printf(PROMPT)

		// Read a line from the input.
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Get the line from the scanner.
		line := scanner.Text()

		// Create a new lexer.
		l := lexer.New(line)

		// Loop through the tokens returned by the lexer.
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			// Print the token type and literal.
			fmt.Printf("%+v\n", tok)
		}
	}
}
