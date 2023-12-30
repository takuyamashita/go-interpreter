package parser

import (
	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/token"
)

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {

	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	exp.Arguments = p.parseCallArguments()

	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {

	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {

		// Read the next token.
		p.nextToken()

		return args
	}

	// ( x , y , z )
	// ^curToken
	p.nextToken()
	// ( x , y , z )
	//   ^curToken

	// Parse the first argument.
	// ( x , y , z )
	//   ^curToken is x so parseExpression(LOWEST) returns x as ast.Identifier
	//
	// if the first argument is statement, parseExpression(LOWEST) returns ast.Statement
	// ( 1 + 2 , y , z )
	//   ^curToken is 1 so parseExpression(LOWEST) returns 1 + 2 as ast.InfixExpression
	args = append(args, p.parseExpression(LOWEST))

	// Parse the rest of the arguments.
	for p.peekTokenIs(token.COMMA) {

		p.nextToken()
		// ( x , y , z )
		//     ^curToken
		p.nextToken()
		// ( x , y , z )
		//       ^curToken

		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}
