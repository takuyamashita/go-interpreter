package parser

import "github.com/takuyamashita/go-interpreter/ast"

func (p *Parser) parsePrefixExpression() ast.Expression {

	// +1, -1, !true, !false, etc.

	// Create a new PrefixExpression AST node.
	expression := &ast.PrefixExpression{
		Token:    p.curToken, // The prefix token, e.g. +, -, !
		Operator: p.curToken.Literal,
	}

	// Read the next token.
	p.nextToken()

	// Parse the right-hand side of the expression.
	expression.Right = p.parseExpression(PREFIX)

	return expression
}
