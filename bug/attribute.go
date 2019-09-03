package bug

import (
)

// An attribute is a named property in a bug (value string only ?)

type Attribute struct {
	name string
	value string
}

func (a Attribute) Name() string {
	if a.name == "" {
		panic("attribute without name")
	}
	return a.name
}

func (a Attribute) Value() string {
	return a.value
}

