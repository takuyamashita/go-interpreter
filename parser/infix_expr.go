package parser

import "github.com/takuyamashita/go-interpreter/ast"

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	// Create a new InfixExpression AST node.
	expression := &ast.InfixExpression{
		Token:    p.curToken, // The operator token, e.g. +, -, *, /, etc.
		Operator: p.curToken.Literal,
		Left:     left,
	}

	// Get the current precedence.
	precedence := p.curPrecedence()

	// Read the next token.
	p.nextToken()

	// Parse the right-hand side of the expression.
	expression.Right = p.parseExpression(precedence)

	return expression
}
