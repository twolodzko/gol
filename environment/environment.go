package environment

import (
	"fmt"

	"github.com/twolodzko/gol/types"
)

type (
	Any    = types.Any
	Bool   = types.Bool
	Int    = types.Int
	Float  = types.Float
	String = types.String
	Symbol = types.Symbol
	List   = types.List
)

type Env struct {
	Objects map[Symbol]Any
	Parent  *Env
}

func NewEnv(env *Env) *Env {
	objs := make(map[Symbol]Any)
	return &Env{Objects: objs, Parent: env}
}

func (env *Env) Find(sym Symbol) (*Env, error) {
	_, ok := env.Objects[sym]

	// recursive search over enviroment
	if !ok {
		if env.Parent != nil {
			return env.Parent.Find(sym)
		} else {
			return nil, fmt.Errorf("unable to resolve %s in this context", sym)
		}
	}

	return env, nil
}

func (env *Env) Get(sym Symbol) (Any, error) {
	env, err := env.Find(sym)
	if err != nil {
		return nil, err
	}
	return env.Objects[sym], nil
}

func (env *Env) Set(sym Symbol, val Any) error {
	env.Objects[sym] = val
	return nil
}
