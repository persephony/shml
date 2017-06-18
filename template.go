package shml

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

var errUnsupportedTransform = errors.New("unsupported transform")

// Template holds the template and parsed variables from the input
type Template struct {
	buf  []byte
	vars *templateVariables
}

// New creates a new Template instance
func New() *Template {
	return &Template{
		vars: &templateVariables{
			order: make([]string, 0),
			m:     make(map[string]*position),
		},
	}
}

// Parse parses the input buffer indexing all variables
func (t *Template) Parse(buf []byte) {
	t.buf = buf

	pos := newPosition()
	for i := 0; i < len(t.buf); i++ {

		switch t.buf[i] {
		case '$':
			if t.buf[i+1] == '{' {
				pos.s = i
				i += 2
			}

		case '\\':
			i++

		case '}':
			if pos.s >= 0 {
				pos.e = i
				// Index by var in template
				val := string(pos.value(t.buf))
				t.vars.set(val, pos)
				// Reset position
				pos = newPosition()
			}

		}

	}

}

// ExecuteIndex applies the data index to the template
func (t *Template) ExecuteIndex(idx ContextVariables) ([]byte, error) {
	if idx == nil {
		return nil, fmt.Errorf("index is nil")
	}

	sort.Sort(t.vars)

	out := []byte{}
	var last int

	err := t.vars.iter(func(rawkey string, pos *position) error {
		// get key name from buffer
		key := string(pos.name(t.buf))
		val, ok := idx[key]
		if !ok {
			return nil
		}
		//log.Println("KEY", key)

		var (
			ival = val.Interface()
			b    []byte
			err  error
		)

		if trns := pos.transform(t.buf); trns != nil {

			if b, err = applyTransform(ival, string(trns)); err != nil {
				return err
			}

		} else {
			b = []byte(fmt.Sprintf("%v", ival))
		}
		// copy data
		to := make([]byte, pos.s-last)
		copy(to, t.buf[last:pos.s])
		// append value
		to = append(to, b...)
		// copy to output buffer
		out = append(out, to...)
		// update position in original buffer
		last = pos.e + 1
		return nil
	})

	if last < len(t.buf)-1 {
		out = append(out, t.buf[last:]...)
	}

	return out, err
}

// Execute generates an index of the data and applies it to the template
func (t *Template) Execute(data interface{}) ([]byte, error) {
	ctx, err := BuildIndex(data)
	if err != nil {
		return nil, err
	}

	return t.ExecuteIndex(ctx)
}

func applyTransform(data interface{}, transform string, args ...string) ([]byte, error) {
	var err error
	switch transform {
	case "json":
		if len(args) == 0 {
			return json.Marshal(data)
		}

		err = fmt.Errorf("invalid transform args: %v", args)

	default:
		err = errUnsupportedTransform
	}

	return nil, err
}
