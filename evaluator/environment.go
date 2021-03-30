package evaluator

import (
	"errors"
	"fmt"

	"github.com/twolodzko/goal/objects"
)

var baseEnv = &Enviroment{buildins, nil}

type Enviroment struct {
	objects map[string]objects.Object
	parent  *Enviroment
}

func (env *Enviroment) Get(name string) (objects.Object, error) {
	val, ok := env.objects[name]
	if !ok {
		if env.parent != nil {
			return env.parent.Get(name)
		} else {
			return nil, fmt.Errorf("object %s not found", name)
		}
	}
	return val, nil
}

func (env *Enviroment) Set(name string, val objects.Object) error {
	if env.parent == nil {
		return errors.New("cannot set values in base enviroment")
	}
	env.objects[name] = val
	return nil
}
