package object

import (
	"bytes"
	"strings"

	"github.com/takuyamashita/go-interpreter/ast"
)

const (
	FUNCTION_OBJ = "FUNCTION"
)

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement

	// local environment
	Env *Environment
}

func (f *Function) Type() ObjectType {

	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {

	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {

		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}
