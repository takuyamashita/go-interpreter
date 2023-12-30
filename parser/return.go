package parser

import (
	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/token"
)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {

	// Create a new ReturnStatement AST node.
	stmt := &ast.ReturnStatement{Token: p.curToken}

	// Read the next token.
	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {

		// Read the next token.
		p.nextToken()
	}

	return stmt
}
