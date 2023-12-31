package object

const (
	BOOLEAN_OBJ ObjectType = "BOOLEAN"
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return b.Inspect()
}
