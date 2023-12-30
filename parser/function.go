package parser

import (
	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/token"
)

func (p *Parser) parseFunctionLiteral() ast.Expression {

	// Create a new FunctionLiteral AST node.
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Check if the next token is "(".
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Parse the function parameters.
	lit.Parameters = p.parseFunctionParameters()

	// Check if the next token is "{".
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// Parse the function body.
	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {

	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		// no parameters

		p.nextToken()

		return identifiers
	}

	// ( x , y , z )
	// ^curToken
	p.nextToken()
	// ( x , y , z )
	//   ^curToken

	// first parameter
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Append the identifier to the identifiers array.
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {

		p.nextToken()
		// ( x , y , z )
		//     ^curToken

		p.nextToken()
		// ( x , y , z )
		//       ^curToken

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}
