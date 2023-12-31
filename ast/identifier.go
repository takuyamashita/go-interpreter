package ast

import "github.com/takuyamashita/go-interpreter/token"

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {

	return i.Token.Literal
}

func (i *Identifier) String() string {

	return i.Value
}
