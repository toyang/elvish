package eval

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/xiaq/persistent/hashmap"
)

// Map is a map from string to Value.
type Map struct {
	inner hashmap.HashMap
}

type HasKeyer interface {
	HasKey(k Value) bool
}

var _ MapLike = Map{}

// NewMap creates a new Map from an inner HashMap.
func NewMap(inner hashmap.HashMap) Map {
	return Map{inner}
}

// ConvertToMap converts a native Go map to Map.
func ConvertToMap(m map[Value]Value) Map {
	inner := hashmap.Empty
	for k, v := range m {
		inner = inner.Assoc(k, v)
	}
	return NewMap(inner)
}

func (Map) Kind() string {
	return "map"
}

func (m Map) Equal(a interface{}) bool {
	return m == a || eqMapLike(m, a)
}

func (m Map) Hash() uint32 {
	return hashMapLike(m)
}

func (m Map) MarshalJSON() ([]byte, error) {
	// TODO(xiaq): Replace with a more efficient implementation.
	mm := map[string]Value{}
	for it := m.inner.Iterator(); it.HasElem(); it.Next() {
		k, v := it.Elem()
		mm[ToString(k.(Value))] = v.(Value)
	}
	return json.Marshal(mm)
}

func (m Map) Repr(indent int) string {
	var builder MapReprBuilder
	builder.Indent = indent
	for it := m.inner.Iterator(); it.HasElem(); it.Next() {
		k, v := it.Elem()
		builder.WritePair(k.(Value).Repr(indent+1), indent+2, v.(Value).Repr(indent+2))
	}
	return builder.String()
}

func (m Map) Len() int {
	return m.inner.Len()
}

func (m Map) IndexOne(idx Value) Value {
	v, ok := m.inner.Get(idx)
	if !ok {
		throw(errors.New("no such key: " + idx.Repr(NoPretty)))
	}
	return v.(Value)
}

func (m Map) Assoc(k, v Value) Value {
	return Map{m.inner.Assoc(k, v)}
}

func (m Map) IterateKey(f func(Value) bool) {
	for it := m.inner.Iterator(); it.HasElem(); it.Next() {
		k, _ := it.Elem()
		if !f(k.(Value)) {
			break
		}
	}
}

func (m Map) IteratePair(f func(Value, Value) bool) {
	for it := m.inner.Iterator(); it.HasElem(); it.Next() {
		k, v := it.Elem()
		if !f(k.(Value), v.(Value)) {
			break
		}
	}
}

func (m Map) HasKey(k Value) bool {
	_, ok := m.inner.Get(k)
	return ok
}

// MapReprBuilder helps building the Repr of a Map. It is also useful for
// implementing other Map-like values. The zero value of a MapReprBuilder is
// ready to use.
type MapReprBuilder struct {
	ListReprBuilder
}

func (b *MapReprBuilder) WritePair(k string, indent int, v string) {
	if indent > 0 {
		b.WriteElem("&" + k + "=\n" + strings.Repeat(" ", indent) + v)
	} else {
		b.WriteElem("&" + k + "=" + v)
	}
}

func (b *MapReprBuilder) String() string {
	if b.buf.Len() == 0 {
		return "[&]"
	}
	return b.ListReprBuilder.String()
}
