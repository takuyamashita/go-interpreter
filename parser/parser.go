package parser

import (
	"fmt"

	"github.com/takuyamashita/go-interpreter/ast"
	"github.com/takuyamashita/go-interpreter/lexer"
	"github.com/takuyamashita/go-interpreter/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	// curToken and peekToken are the current token and the next token.
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {

	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Register prefix parse functions.
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

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

func (p *Parser) parseIdentifier() ast.Expression {

	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {

	// Check if the current token type has a prefix parse function.
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {

		// If it does not, add an error to the parser.
		return nil
	}

	// Call the prefix parse function.
	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {

	// Create a new ExpressionStatement AST node.
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	// Parse the expression.
	stmt.Expression = p.parseExpression(LOWEST)

	// Check if the next token is a ";".
	if p.peekTokenIs(token.SEMICOLON) {

		// Read the next token.
		p.nextToken()
	}

	return stmt
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

		// Read the next token.
		p.nextToken()
	}

	return stmt

}

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
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {

	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {

	// Create an error message.
	msg := "expected next token to be %s, got %s instead"
	msg = fmt.Sprintf(msg, t, p.peekToken.Type)

	// Add the error to the parser.
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {

	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {

	p.infixParseFns[tokenType] = fn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)
