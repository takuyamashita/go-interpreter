package object

import "fmt"

const (
	INTEGER_OBJ ObjectType = "INTEGER"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
