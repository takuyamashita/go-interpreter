package parser

import (
	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/token"
)

func (p *Parser) parseIfExpression() ast.Expression {

	// Create a new IfExpression AST node.
	expression := &ast.IfExpression{Token: p.curToken}

	// Check if the next token is "(".
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Read the next token.
	p.nextToken()

	// Parse the expression.
	expression.Condition = p.parseExpression(LOWEST)

	// Check if the next token is ")".
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Check if the next token is "{".
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// Parse the consequence.
	expression.Consequence = p.parseBlockStatement()

	// Check if the next token is "else".
	if p.peekTokenIs(token.ELSE) {

		// Read the next token.
		p.nextToken()

		// Check if the next token is "{".
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		// Parse the alternative.
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}
