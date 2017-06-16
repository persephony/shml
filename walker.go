package shml

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/reflectwalk"
)

// ReflectWalker walks a data structure constructing the context variables
type ReflectWalker struct {
	vars   ContextVariables
	prefix string
}

// NewReflectWalker creates a new instance of ReflectWalker
func NewReflectWalker() *ReflectWalker {
	return &ReflectWalker{vars: make(ContextVariables)}
}

func (t *ReflectWalker) levelupPrefix(name string) {
	if t.prefix == "" {
		t.prefix = name
	} else {
		t.prefix += "." + name
	}
}

func (t *ReflectWalker) leveldownPrefix() {
	arr := strings.Split(t.prefix, ".")
	if len(arr) == 1 {
		t.prefix = ""
	} else {
		t.prefix = strings.Join(arr[:len(arr)-1], ".")
	}
}

func (t *ReflectWalker) path(name string) string {
	if t.prefix == "" {
		return name
	}
	return t.prefix + "." + name
}

// Enter is required for reflectwalk interface in order to use Exit
func (t *ReflectWalker) Enter(l reflectwalk.Location) error {
	return nil
}

// Exit trims the last part of the prefix for data structures only
func (t *ReflectWalker) Exit(l reflectwalk.Location) error {
	switch l {
	case reflectwalk.Struct, reflectwalk.Map:
		t.leveldownPrefix()
	}
	return nil
}

// Struct is required for reflectwalk interface
func (t *ReflectWalker) Struct(v reflect.Value) error {
	return nil
}

// StructField sets the prefix if the value is a nestable data structure and updates the index
func (t *ReflectWalker) StructField(sf reflect.StructField, v reflect.Value) error {
	key := t.path(sf.Name)
	t.vars[key] = v

	kind := v.Kind()
	switch kind {
	case reflect.Struct, reflect.Map:
		t.levelupPrefix(sf.Name)

	case reflect.Interface:
		if v.IsNil() {
			break
		}
		knd := reflect.TypeOf(v.Interface()).Kind()
		switch knd {
		case reflect.Map, reflect.Struct:
			t.levelupPrefix(sf.Name)

		}

	}

	return nil
}

// Map is required for reflectwalk interface in order to use MapElem
func (t *ReflectWalker) Map(m reflect.Value) error {
	return nil
}

// MapElem converts a they key to a string.  If the value is a data structure prefix is also updated
func (t *ReflectWalker) MapElem(m, k, v reflect.Value) error {

	var (
		keyKind = k.Kind()
		keystr  string
	)

	switch keyKind {
	case reflect.String:
		keystr = k.Interface().(string)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := k.Interface()
		keystr = fmt.Sprintf("%d", i)
	}

	t.vars[t.path(keystr)] = v

	valKind := v.Kind()
	switch valKind {
	case reflect.Map, reflect.Struct:
		t.levelupPrefix(keystr)

	case reflect.Interface:
		if v.IsNil() {
			break
		}
		knd := reflect.TypeOf(v.Interface()).Kind()
		switch knd {
		case reflect.Map, reflect.Struct:
			t.levelupPrefix(keystr)
		}

	}

	return nil
}
