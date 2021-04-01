package evaluator

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

var baseEnv = &Enviroment{buildins, nil}

type Enviroment struct {
	objects map[string]Any
	parent  *Enviroment
}

func (env *Enviroment) Get(sym Symbol) (Any, error) {
	name := string(sym)
	val, ok := env.objects[name]
	if !ok {
		if env.parent != nil {
			return env.parent.Get(sym)
		} else {
			return nil, fmt.Errorf("object %s not found", name)
		}
	}
	return val, nil
}

func (env *Enviroment) Set(sym Symbol, val Any) error {
	name := string(sym)
	env.objects[name] = val
	return nil
}
