package shml

import (
	"encoding/json"
	"testing"

	"github.com/d3sw/floop/types"
	"github.com/mitchellh/reflectwalk"
)

var (
	testTemplate = `Event:
  Type: ${Type}
  Meta:
    ${Meta.int-key}
    ${Meta.map-string-key.float-key}
    ${Meta.kstring}/${Meta.map-string-key.string-key}

  Data: ${Data|json}

\${escaped}
-----------`
	testTemplateOut = `Event:
  Type: begin
  Meta:
    1
    1.2
    string/value

  Data: {"Type":"http","Transform":null,"Context":null,"Config":{"key":"value"},"IgnoreErrors":false}

\${escaped}
-----------`
	testEvent = &types.Event{
		Type: types.EventTypeBegin,
		Meta: map[string]interface{}{
			"kstring": "string",
			"int-key": 1,
			"map-string-key": map[string]interface{}{
				"string-key": "value",
				"float-key":  1.2,
			},
			"map-int-key": map[int]string{
				2: "two",
				5: "five",
			},
			"map-uint-key": map[uint]string{
				uint(9): "nine",
			},
			"list-key": []interface{}{"foo", 2},
		},
		Data: types.HandlerConfig{
			Type: "http",
			Options: map[string]interface{}{
				"key": "value",
			},
		},
	}
)

func Test_ReflectWalker(t *testing.T) {

	walker := NewReflectWalker()
	if err := reflectwalk.Walk(testEvent, walker); err != nil {
		t.Fatal(err)
	}

	vars := walker.vars
	if len(vars) < 19 {
		t.Fatal("not all keys present")
	}

	b, _ := json.MarshalIndent(vars.keys(), "", "  ")
	t.Logf("%s\n", b)
}

func Test_Template(t *testing.T) {
	tmpl := New()
	tmpl.Parse([]byte(testTemplate))

	// for k, v := range tmpl.vars {
	// 	t.Log(k, v)
	// }

	out, err := tmpl.Execute(testEvent)
	if err != nil {
		t.Fatal(err)
	}

	// if string(out) != testTemplateOut {
	// 	t.Fatal("failed to execute template")
	// }

	t.Logf("%s", out)
}
