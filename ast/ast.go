package ast

import (
	"bytes"

	"github.com/takuyamashita/go-interpreter/token"
)

// let x = 10;
// let add = fn(x, y) {
// 	x + y;
// };
// the above code is aggregration of the statements below:
// let <identifier> = <expression>;

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {

	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {

	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// ex ) 5 + 5;
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {

	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {

	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {

	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {

	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// ex ) 5 + 3;
type InfixExpression struct {
	Token    token.Token // the operator token, e.g. +
	Left     Expression  // e.g. 5
	Operator string      // e.g. +
	Right    Expression  // e.g. 3
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {

	return ie.Token.Literal
}

func (ie *InfixExpression) String() string {

	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
