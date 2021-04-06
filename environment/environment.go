package environment

import (
	"fmt"

	. "github.com/twolodzko/gol/types"
)

type Env struct {
	Objects map[Symbol]Any
	Parent  *Env
}

func NewEnv(env *Env) *Env {
	objs := make(map[Symbol]Any)
	return &Env{Objects: objs, Parent: env}
}

func (env *Env) Get(sym Symbol) (Any, error) {
	val, ok := env.Objects[sym]

	// recursive search over enviroment
	if !ok {
		if env.Parent != nil {
			return env.Parent.Get(sym)
		} else {
			return nil, fmt.Errorf("unable to resolve %s in this context", sym)
		}
	}

	return val, nil
}

func (env *Env) Set(sym Symbol, val Any) error {
	env.Objects[sym] = val
	return nil
}
