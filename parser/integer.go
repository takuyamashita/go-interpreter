package parser

import (
	"fmt"
	"strconv"

	"github.com/takuyamashita/go-interpreter/ast"
)

func (p *Parser) parseIntegerLiteral() ast.Expression {

	// Create a new IntegerLiteral AST node.
	lit := &ast.IntegerLiteral{Token: p.curToken}

	// Try to parse the literal as an integer.
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {

		// If it fails, add an error to the parser.
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)

		return nil
	}

	// Set the literal value.
	lit.Value = value

	return lit
}
