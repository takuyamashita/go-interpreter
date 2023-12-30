package parser

import (
	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/token"
)

func (p *Parser) parseLetStatement() *ast.LetStatement {

	// Create a new LetStatement AST node.
	stmt := &ast.LetStatement{Token: p.curToken}

	// Check if the next token is an identifier.
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Set the identifier name.
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Check if the next token is an "=".
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {

		// Read the next token.
		p.nextToken()
	}

	return stmt

}
