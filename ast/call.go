package ast

import (
	"bytes"
	"strings"

	"github.com/takuyamashita/go-interpreter/token"
)

type CallExpression struct {

	// The '(' token.
	Token token.Token

	// Identifier or FunctionLiteral
	// e.g. add(2, 3) -> add is Identifier
	// e.g. fn(x, y) { x + y; }(2, 3) -> fn(x, y) { x + y; } is function literal
	Function Expression

	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {

	return ce.Token.Literal
}

func (ce *CallExpression) String() string {

	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {

		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
