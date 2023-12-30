package parser

import "github.com/takuyamashita/go-interpreter/ast"

func (p *Parser) parseIdentifier() ast.Expression {

	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}
