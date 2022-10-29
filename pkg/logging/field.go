package logging

type Field struct {
	Label string
	Value interface{}
}

func NewField(label string, value interface{}) Field {
	return Field{
		Label: label,
		Value: value,
	}
}
