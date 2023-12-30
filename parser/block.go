package parser

import (
	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/token"
)

func (p *Parser) parseBlockStatement() *ast.BlockStatement {

	// Create a new BlockStatement AST node.
	block := &ast.BlockStatement{Token: p.curToken}

	// Initialize the statements array.
	block.Statements = []ast.Statement{}

	// Read the next token.
	p.nextToken()

	// Parse the statements until we encounter a "}" token.
	for !p.curTokenIs(token.RBRACE) {

		// Parse the statement.
		stmt := p.parseStatement()

		// Append the statement to the statements array.
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		// Read the next token.
		p.nextToken()
	}

	return block
}
