package parser

import "github.com/takuyamashita/go-interpreter/ast"

func (p *Parser) parseStringLiteral() ast.Expression {

	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}
