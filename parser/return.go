package parser

import (
	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/token"
)

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {

	// return 5 + 10;
	//  ^curToken

	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	// return 5 + 10;
	//        ^curToken

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {

		p.nextToken()
		// return 5 + 10;
		//              ^curToken

		return stmt
	}

	return stmt
}
