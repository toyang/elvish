package eval

import (
	"errors"
	"fmt"

	"github.com/elves/elvish/parse"
)

var (
	ErrIndexMustBeString = errors.New("index must be string")
)

// Struct is like a Map with fixed keys.
type Struct struct {
	Descriptor *StructDescriptor
	Fields     []Value
}

var (
	_ Value   = (*Struct)(nil)
	_ MapLike = (*Struct)(nil)
)

func (*Struct) Kind() string {
	return "map"
}

// Equal returns true if the rhs is MapLike and all pairs are equal.
func (s *Struct) Equal(rhs interface{}) bool {
	return s == rhs || eqMapLike(s, rhs)
}

func (s *Struct) Hash() uint32 {
	return hashMapLike(s)
}

func (s *Struct) Repr(indent int) string {
	var builder MapReprBuilder
	builder.Indent = indent
	for i, name := range s.Descriptor.fieldNames {
		builder.WritePair(parse.Quote(name), indent+2, s.Fields[i].Repr(indent+2))
	}
	return builder.String()
}

func (s *Struct) Len() int {
	return len(s.Descriptor.fieldNames)
}

func (s *Struct) IndexOne(idx Value) Value {
	return s.Fields[s.index(idx)]
}

func (s *Struct) Assoc(k, v Value) Value {
	i := s.index(k)
	fields := make([]Value, len(s.Fields))
	copy(fields, s.Fields)
	fields[i] = v
	return &Struct{s.Descriptor, fields}
}

func (s *Struct) IterateKey(f func(Value) bool) {
	for _, field := range s.Descriptor.fieldNames {
		if !f(String(field)) {
			break
		}
	}
}

func (s *Struct) IteratePair(f func(Value, Value) bool) {
	for i, field := range s.Descriptor.fieldNames {
		if !f(String(field), s.Fields[i]) {
			break
		}
	}
}

func (s *Struct) HasKey(k Value) bool {
	index, ok := k.(String)
	if !ok {
		return false
	}
	_, ok = s.Descriptor.fieldIndex[string(index)]
	return ok
}

func (s *Struct) index(idx Value) int {
	index, ok := idx.(String)
	if !ok {
		throw(ErrIndexMustBeString)
	}
	i, ok := s.Descriptor.fieldIndex[string(index)]
	if !ok {
		throw(fmt.Errorf("no such field: %s", index.Repr(NoPretty)))
	}
	return i
}

type StructDescriptor struct {
	fieldNames []string
	fieldIndex map[string]int
}

func NewStructDescriptor(fields ...string) *StructDescriptor {
	fieldNames := append([]string(nil), fields...)
	fieldIndex := make(map[string]int)
	for i, name := range fieldNames {
		fieldIndex[name] = i
	}
	return &StructDescriptor{fieldNames, fieldIndex}
}
