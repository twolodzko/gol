package environment

import "testing"

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

	env := NewEnv(nil)

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

	newEnv := NewEnv(env)

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

	// find variable in top env
	foundEnv, err := newEnv.Find(Symbol("w"))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// it contains the searched object
	if _, ok := foundEnv.Objects[Symbol("w")]; !ok {
		t.Errorf("did not find the correct env: %v", foundEnv)
	}
	// it doesn't contain the object from parent env
	if _, ok := foundEnv.Objects[Symbol("x")]; ok {
		t.Errorf("did not find the correct env: %v", foundEnv)
	}

	// find object in parent env
	foundEnv, err = newEnv.Find(Symbol("x"))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if _, ok := foundEnv.Objects[Symbol("x")]; !ok {
		t.Errorf("did not find the correct env: %v", foundEnv)
	}
	if _, ok := foundEnv.Objects[Symbol("y")]; !ok {
		t.Errorf("did not find the correct env: %v", foundEnv)
	}
	// it doesn't contain the object from child env
	if _, ok := foundEnv.Objects[Symbol("w")]; ok {
		t.Errorf("did not find the correct env: %v", foundEnv)
	}

}
