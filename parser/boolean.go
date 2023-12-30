package parser

import (
	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/token"
)

func (p *Parser) parseBoolean() ast.Expression {

	// Create a new Boolean AST node.
	expression := &ast.Boolean{Token: p.curToken}

	// Set the boolean value.
	expression.Value = p.curTokenIs(token.TRUE)

	return expression
}
