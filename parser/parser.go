package parser

import (
	"fmt"
	"strconv"

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

var (
	precedences = map[token.TokenType]int{
		token.EQ:     EQUALS,
		token.NOT_EQ: EQUALS,
		token.LT:     LESSGREATER,
		token.GT:     LESSGREATER,
		token.PLUS:   SUM,
		token.MINUS:  SUM,
		token.SLASH:  PRODUCT,
		token.ASTER:  PRODUCT,
		token.MOD:    PRODUCT,
	}
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
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	// Rgister infix parse functions.
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTER, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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

func (p *Parser) peekPrecedence() int {

	// Check if the precedence exists in the map.
	if p, ok := precedences[p.peekToken.Type]; ok {

		// If it does, return it.
		return p
	}

	// If it does not, return LOWEST.
	return LOWEST
}

func (p *Parser) curPrecedence() int {

	// Check if the precedence exists in the map.
	if p, ok := precedences[p.curToken.Type]; ok {

		// If it does, return it.
		return p
	}

	// If it does not, return LOWEST.
	return LOWEST
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

		p.noPrefixParseFnError(p.curToken.Type)

		return nil
	}

	// Call the prefix parse function.
	leftExp := prefix()

	// Loop through the tokens until we reach a semicolon or the precedence is
	// lower than the next precedence.
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {

		// Check if the next token type has an infix parse function.
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {

			// If it does not, return the left-hand side expression.
			return leftExp
		}

		// Read the next token.
		p.nextToken()

		// Call the infix parse function.
		leftExp = infix(leftExp)
	}

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

func (p *Parser) parsePrefixExpression() ast.Expression {

	// +1, -1, !true, !false, etc.

	// Create a new PrefixExpression AST node.
	expression := &ast.PrefixExpression{
		Token:    p.curToken, // The prefix token, e.g. +, -, !
		Operator: p.curToken.Literal,
	}

	// Read the next token.
	p.nextToken()

	// Parse the right-hand side of the expression.
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	// 5 + 5, 5 - 5, etc.

	// Create a new InfixExpression AST node.
	expression := &ast.InfixExpression{
		Token:    p.curToken, // The operator token, e.g. +, -, *, /, etc.
		Operator: p.curToken.Literal,
		Left:     left,
	}

	// Get the current precedence.
	precedence := p.curPrecedence()

	// Read the next token.
	p.nextToken()

	// Parse the right-hand side of the expression.
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseIntegerLiteral() ast.Expression {

	// Create a new IntegerLiteral AST node.
	lit := &ast.IntegerLiteral{Token: p.curToken}

	// Try to parse the literal as an integer.
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {

		// If it fails, add an error to the parser.
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)

		return nil
	}

	// Set the literal value.
	lit.Value = value

	return lit
}

func (p *Parser) parseBoolean() ast.Expression {

	// Create a new Boolean AST node.
	expression := &ast.Boolean{Token: p.curToken}

	// Set the boolean value.
	expression.Value = p.curTokenIs(token.TRUE)

	return expression
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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {

	msg := fmt.Sprintf("no prefix parse function for %s found", t)

	p.errors = append(p.errors, msg)
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
