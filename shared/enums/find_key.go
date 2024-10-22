package enums

type FieldOperator int

const (
	Equals FieldOperator = iota
	GreateThan
	GreaterThanOrEqual
	In
	LessThan
	LessThanOrEqual
	NotEqual
	NotIn
	Expression
	Modulo
	Regex
	textStartsWith
	TextContains
	TextEndsWith
	All
	Lenght
)

type logicalOperator int

const (
	And logicalOperator = iota
	Or
	AndNot
	OrNot
)
