package lexer

import "github.com/takuyamashita/go-interpreter/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New returns a new Lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	// Check if we reached the end of the input string.
	if l.readPosition >= len(l.input) {
		// If we reached the end of the input string, set ch to 0, which is the
		// ASCII code for the "NUL" character and is also the value of l.ch when
		// we reach the end of the input string.
		l.ch = 0
	} else {
		// If we haven't reached the end of the input string, set l.ch to the
		// ASCII code of the character at the current read position.
		l.ch = l.input[l.readPosition]
	}

	// Increment the position and readPosition by 1.
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {

	var tok token.Token

	// Skip whitespace characters.
	l.skipWhitespace()

	switch l.ch {
	// Check if l.ch is a letter.
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
		// Check if the next character is '='.
		/*
			if l.peekChar() == '=' {
				// If the next character is '=', we have a two-character token.
				// Read the next character.
				ch := l.ch
				l.readChar()

				// Set the token type to EQ and the literal to "=="
				tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
			} else {
				// If the next character is not '=', we have a one-character token.
				// Set the token type to ASSIGN and the literal to "=".
				tok = newToken(token.ASSIGN, l.ch)
			}
		*/
	case '+':
		// Set the token type to PLUS and the literal to "+".
		tok = newToken(token.PLUS, l.ch)
	case '-':
		// Set the token type to MINUS and the literal to "-".
		tok = newToken(token.MINUS, l.ch)
	/*
		case '!':
			// Check if the next character is '='.
			if l.peekChar() == '=' {
				// If the next character is '=', we have a two-character token.
				// Read the next character.
				ch := l.ch
				l.readChar()

				// Set the token type to NE and the literal to "!="
				tok = token.Token{Type: token.NE, Literal: string(ch) + string(l.ch)}
			} else {
				// If the next character is not '=', we have a one-character token.
				// Set the token type to BANG and the literal to "!".
				tok = newToken(token.BANG, l.ch)
			}
	*/
	case '!':
		// Set the token type to BANG and the literal to "!".
		tok = newToken(token.BANG, l.ch)
	case '*':
		// Set the token type to ASTER and the literal to "*".
		tok = newToken(token.ASTER, l.ch)
	case '/':
		// Set the token type to SLASH and the literal to "/".
		tok = newToken(token.SLASH, l.ch)
	case '%':
		// Set the token type to MOD and the literal to "%".
		tok = newToken(token.MOD, l.ch)

	case '<':
		// Set the token type to LT and the literal to "<".
		tok = newToken(token.LT, l.ch)
	case '>':
		// Set the token type to GT and the literal to ">".
		tok = newToken(token.GT, l.ch)
	case ',':
		// Set the token type to COMMA and the literal to ",".
		tok = newToken(token.COMMA, l.ch)
	case ';':
		// Set the token type to SEMICOLON and the literal to ";".
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		// Set the token type to COLON and the literal to ":".
		tok = newToken(token.COLON, l.ch)
	case '.':
		// Set the token type to PERIOD and the literal to ".".
		tok = newToken(token.PERIOD, l.ch)
	case '|':
		// Set the token type to PIPE and the literal to "|".
		tok = newToken(token.PIPE, l.ch)
	case '(':
		// Set the token type to LPAREN and the literal to "(".
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		// Set the token type to RPAREN and the literal to ")".
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		// Set the token type to LBRACE and the literal to "{".
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		// Set the token type to RBRACE and the literal to "}".
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		// Set the token type to LBRACK and the literal to "[".
		tok = newToken(token.LBRACK, l.ch)
	case ']':
		// Set the token type to RBRACK and the literal to "]".
		tok = newToken(token.RBRACK, l.ch)
	/*
		case '"':
			// Set the token type to STRING and the literal to the string literal.
			tok.Type = token.STRING
			tok.Literal = l.readString()
	*/
	case 0:
		// Set the token type to EOF and the literal to "".
		tok.Literal = ""
		tok.Type = token.EOF
	default:

		if isLetter(l.ch) {

			// Set the token type to IDENT and the literal to the identifier.
			// ex) let, add, foobar, x, y, ...
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)

			return tok

		} else if isDigit(l.ch) {

			// Set the token type to INT and the literal to the integer.
			tok.Literal = l.readNumber()
			tok.Type = token.INT

			return tok

		} else {

			// Set the token type to ILLEGAL and the literal to the illegal character.
			tok = newToken(token.ILLEGAL, l.ch)
		}

	}

	l.readChar()

	return tok
}

func (l *Lexer) readIdentifier() string {

	// Save the current position.
	position := l.position

	// Read characters until we encounter a non-letter-character.
	for isLetter(l.ch) {
		l.readChar()
	}

	// Return the string from the saved position to the current position.
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {

	// Save the current position.
	position := l.position

	// Read characters until we encounter a non-digit-character.
	for isDigit(l.ch) {
		l.readChar()
	}

	// Return the string from the saved position to the current position.
	return l.input[position:l.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	// Return a Token with the given type and literal.
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	// Check if ch is a letter.
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	// Check if ch is a digit.
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {

	// Read characters until we encounter a non-whitespace-character.
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {

		l.readChar()
	}
}
