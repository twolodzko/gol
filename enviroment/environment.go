package enviroment

import (
	"fmt"

	. "github.com/twolodzko/goal/types"
)

type Env struct {
	Objects map[Symbol]Any
	Parent  *Env
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

	// resolve symbols pointing to other symbols
	switch val := val.(type) {
	case Symbol:
		return env.Get(val)
	default:
		return val, nil
	}
}

func (env *Env) Set(sym Symbol, val Any) error {
	env.Objects[sym] = val
	return nil
}
