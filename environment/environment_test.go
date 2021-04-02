package environment

import (
	"testing"

	. "github.com/twolodzko/goal/types"
)

func TestEnv(t *testing.T) {
	var (
		err    error
		result Any
	)

	var testCases = []struct {
		name  Symbol
		value Any
	}{
		{Symbol("x"), Int(1)},
		{Symbol("y"), Int(2)},
		{Symbol("z"), Symbol("x")},
		{Symbol("x"), Int(3)},
	}

	env := NewEnv()

	for _, tt := range testCases {
		err = env.Set(tt.name, tt.value)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		result, err = env.Get(tt.name)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if tt.value != result {
			t.Errorf("expected: %v, got: %v", tt.value, result)
		}
	}

	newEnv := NewEnclosedEnv(env)

	name := Symbol("w")
	value := Int(4)

	err = newEnv.Set(name, value)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	result, err = newEnv.Get(name)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if result != value {
		t.Errorf("expected: %v, got: %v", value, result)
	}

	expected := Int(3)

	result, err = newEnv.Get("x")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if result != expected {
		t.Errorf("expected: %v, got: %v", value, expected)
	}

}
