package types

type ValidationError struct {
	Field string
	Tag   string
	Value string
}
