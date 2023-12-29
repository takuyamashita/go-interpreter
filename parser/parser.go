package parser

import (
	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/lexer"
	"github.com/takuyamashita/go-interpreter/token"
)

type Parser struct {
	l *lexer.Lexer

	// curToken and peekToken are the current token and the next token.
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {

	p := &Parser{l: l}

	// Read two tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {

	// Set curToken to peekToken.
	p.curToken = p.peekToken

	// Read the next token.
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {

	// Create a new Program AST node.
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// Loop through all the tokens until we reach the end of the input.
	for p.curToken.Type != token.EOF {

		// Parse the current statement.
		stmt := p.parseStatement()
		if stmt != nil {

			// Add the statement to the program's statements.
			program.Statements = append(program.Statements, stmt)
		}

		// Read the next token.
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {

	// Check the current token type and call the appropriate function to parse
	// the statement.
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

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

		// End of the statement.
		// Read the next token.
		p.nextToken()
	}

	return stmt

}

func (p *Parser) curTokenIs(t token.TokenType) bool {

	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {

	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {

	// Check if the next token is of the expected type.
	if p.peekTokenIs(t) {

		// If it is, read the next token.
		p.nextToken()
		return true
	} else {

		// If it is not, add an error to the parser.
		//p.peekError(t)
		return false
	}
}
