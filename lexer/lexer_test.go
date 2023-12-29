package lexer

import (
	"testing"

	"github.com/takuyamashita/go-interpreter/token"
)

func TestNextToken(t *testing.T) {

	input := `let five = 5;
let ten = 10;
	
let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// let five = 5;
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		// let ten = 10;
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		// let add = fn(x, y) {
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		// fn(x, y) {
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		// x, y
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		// )
		{token.RPAREN, ")"},
		// {
		{token.LBRACE, "{"},
		// x + y;
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		// }
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		// let result = add(five, ten);
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		// add(five, ten);
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		// five, ten
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		// )
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		// EOF
		{token.EOF, ""},
	}

	// Create a new lexer with the input string.
	l := New(input)

	for i, tt := range tests {
		// Get the next token from the lexer.
		tok := l.NextToken()

		// Check the type of the token.
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		// Check the literal of the token.
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}

}