package evaluator

import (
	. "github.com/twolodzko/goal/types"
)

type Buildin = func([]Any) (Any, error)

// BaseEnv := enviroment.NewEnv()

var Buildins = map[Symbol]Any{
	// type conversions
	// "str":   vectorize(ToString),
	// "int":   vectorize(ToInt),
	// "float": vectorize(ToFloat),
	// type checking
	// "bool?":  vectorize(IsBool),
	// "int?":   vectorize(IsInt),
	// "float?": vectorize(IsFloat),
	// "str?":   vectorize(IsString),
	// "atom?":  vectorize(IsAtom),
	// "list?":  vectorize(IsList),
	// math
	// "%":    FloatMod,
	// "int+": IntSum,
	// "int-": IntDif,
	// "int*": IntMul,
	// "int/": IntDiv,
	// "int%": IntMod,
	// "pow":  Pow,
	// "rem":  Rem,
}

// func vectorize(fn func(Any) (Any, error)) Buildin {
// 	return func(objs []Any) (Any, error) {
// 		if len(objs) == 1 {
// 			return fn(objs[0])
// 		}

// 		var out List
// 		for _, x := range objs {
// 			result, err := fn(x)
// 			if err != nil {
// 				return out, err
// 			}
// 			out = append(out, result)
// 		}
// 		return out, nil
// 	}
// }
