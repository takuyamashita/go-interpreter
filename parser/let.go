package parser

import (
	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/token"
)

func (p *Parser) parseLetStatement() *ast.LetStatement {

	// let x = 5;
	//  ^curToken

	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// let x = 5;
	//     ^curToken

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// let x = 5;
	//       ^curToken

	p.nextToken()
	// let x = 5;
	//         ^curToken

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {

		p.nextToken()
	}

	return stmt

}
