package evaluator

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

var BaseEnv = &Env{buildIns, nil}

type Env struct {
	objects map[Symbol]Any
	parent  *Env
}

func (env *Env) Get(sym Symbol) (Any, error) {
	val, ok := env.objects[sym]
	if !ok {
		if env.parent != nil {
			return env.parent.Get(sym)
		} else {
			return nil, fmt.Errorf("unable to resolve %s in this context", sym)
		}
	}
	return val, nil
}

func (env *Env) Set(sym Symbol, val Any) error {
	env.objects[sym] = val
	return nil
}
