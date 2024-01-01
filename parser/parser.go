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
		token.LPAREN: CALL,
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

	// 前置構文解析関数のマップを初期化する。
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// e.g foo
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	// e.g 5
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	// e.g !, -
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	// e.g true, false
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	// e.g (5 + 5)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	// e.g if (x < y) { x } else { y }
	p.registerPrefix(token.IF, p.parseIfExpression)
	// e.g fn(x, y) { x + y; }
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	// e.g "foobar"
	p.registerPrefix(token.STRING, p.parseStringLiteral)

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
	p.registerInfix(token.LPAREN, p.parseCallExpression)

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

func (p *Parser) ParseProgram() *ast.Program {

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {

		stmt := p.parseStatement()

		if stmt != nil {

			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {

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

	// prefix is the prefix parse function for the current token type.
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {

		p.noPrefixParseFnError(p.curToken.Type)

		return nil
	}

	//
	// e.g 5 + 3 * 2;
	//
	// [first loop]
	// 	leftExp = 5
	// 	precedence = 1(LOWEST)
	// 	peekPrecedence = 4(SUM)
	// 	p.curToken = 5
	//  p.peekToken = +
	//
	//  	nextToken()
	// 		p.curToken = +
	//  	p.peekToken = 3
	//
	//  	parseInfixExpression(5)
	//
	// 			curPrecedence = 4(SUM)
	//
	//         	nextToken()
	// 			p.curToken = 3
	//  		p.peekToken = *
	//
	//  		parseExpression(4)
	//
	//              [first loop]
	//
	//			    leftExp = 3
	// 				precedence = 4(SUM)
	// 				peekPrecedence = 5(PRODUCT)
	// 				p.curToken = 3
	//  			p.peekToken = *
	//
	//              nextToken()
	// 				p.curToken = *
	//  			p.peekToken = 2
	//
	//              parseInfixExpression(3)
	//
	//					curPrecedence = 5(PRODUCT)
	//
	//                	nextToken()
	// 					p.curToken = 2
	//  				p.peekToken = ;
	//
	//  				parseExpression(5)
	//
	//						leftExp = 2
	// 						retrun 2 cuz semi colon
	//

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {

		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {

			return leftExp
		}

		p.nextToken()

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

func (p *Parser) parseGroupedExpression() ast.Expression {

	// Read the next token.
	p.nextToken()

	// Parse the expression.
	expression := p.parseExpression(LOWEST)

	// Check if the next token is ")".
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (p *Parser) curTokenIs(t token.TokenType) bool {

	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {

	return p.peekToken.Type == t
}

// expectPeek checks if the next token in the input (the "peek" token) is of the expected type.
// If the peek token is of the expected type, it advances the parser to the next token and returns true.
// If the peek token is not what was expected, it adds an error to the parser's error list and returns false.
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
