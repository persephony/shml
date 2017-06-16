package shml

import (
	"reflect"
	"sort"

	"github.com/mitchellh/reflectwalk"
)

//
// Syntax:
//   Variables: ${data.key}
//   Functions: $(data|json}
//

// position is the position of a variable in a template
type position struct {
	s int
	e int
}

func newPosition() *position {
	return &position{s: -1, e: -1}
}

// value retuns the value between the markers
func (pos *position) value(buf []byte) []byte {
	return buf[pos.s+2 : pos.e]
}

// name returns the name portion of the variable
func (pos *position) name(buf []byte) []byte {
	val := pos.value(buf)
	for i, b := range val {
		if b == '|' {
			return val[:i]
		}
	}
	return val
}

// transform returns the transform portion of the variable
func (pos *position) transform(buf []byte) []byte {
	val := pos.value(buf)
	for i, b := range val {
		if b == '|' {
			return val[i+1:]
		}
	}
	return nil
}

// templateVariables is a mapping of variable names to their position in the template
type templateVariables struct {
	m     map[string]*position
	order []string
}

func (vars *templateVariables) set(key string, value *position) {
	if _, ok := vars.m[key]; !ok {
		vars.m[key] = value
		vars.order = append(vars.order, key)
	}
}

func (vars *templateVariables) iter(cb func(k string, v *position) error) error {
	for _, key := range vars.order {
		if err := cb(key, vars.m[key]); err != nil {
			return err
		}
	}
	return nil
}

func (vars *templateVariables) Len() int { return len(vars.m) }
func (vars *templateVariables) Swap(i, j int) {
	vars.order[i], vars.order[j] = vars.order[j], vars.order[i]
}
func (vars *templateVariables) Less(i, j int) bool {
	return vars.m[vars.order[i]].s < vars.m[vars.order[j]].s
}

// BuildIndex builds an index for the data structure generating a key to value mapping.
func BuildIndex(data interface{}) (ContextVariables, error) {
	walker := NewReflectWalker()
	if err := reflectwalk.Walk(data, walker); err != nil {
		return nil, err
	}
	return walker.vars, nil
}

// ContextVariables is a mapping of string variable path to actual value
type ContextVariables map[string]reflect.Value

func (vars ContextVariables) keys() []string {
	out := make([]string, len(vars))
	i := 0
	for k := range vars {
		out[i] = k
		i++
	}
	sort.Strings(out)
	return out
}
