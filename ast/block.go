package ast

import (
	"bytes"

	"github.com/takuyamashita/go-interpreter/token"
)

type BlockStatement struct {
	Token      token.Token // The '{' token.
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// Implement the Node interface.
func (bs *BlockStatement) TokenLiteral() string {

	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {

	var out bytes.Buffer

	for _, s := range bs.Statements {

		out.WriteString(s.String())
	}

	return out.String()
}
