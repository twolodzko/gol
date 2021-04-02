package evaluator

import (
	"github.com/twolodzko/goal/enviroment"
	. "github.com/twolodzko/goal/types"
)

type Buildin = func([]Any) (Any, error)

var BaseEnv = &enviroment.Env{Buildins, nil}

var Buildins = map[Symbol]Any{
	"list":    ToList,
	"size":    vectorize(Size),
	"head":    Head,
	"tail":    Tail,
	"nil?":    vectorize(IsNil),
	"error":   Error,
	"eq?":     AreSame,
	"print":   Print,
	"println": PrintLn,
	// type conversions
	"str":   vectorize(ToString),
	"int":   vectorize(ToInt),
	"float": vectorize(ToFloat),
	// type checking
	"bool?":  vectorize(IsBool),
	"int?":   vectorize(IsInt),
	"float?": vectorize(IsFloat),
	"str?":   vectorize(IsString),
	"atom?":  vectorize(IsAtom),
	"list?":  vectorize(IsList),
	// logic
	"true?": vectorize(IsTrue),
	"not":   vectorize(Not),
	"and":   And,
	"or":    Or,
	// math
	"+":    FloatSum,
	"-":    FloatDif,
	"*":    FloatMul,
	"/":    FloatDiv,
	"%":    FloatMod,
	"int+": IntSum,
	"int-": IntDif,
	"int*": IntMul,
	"int/": IntDiv,
	"int%": IntMod,
	"pow":  Pow,
	"rem":  Rem,
}

func vectorize(fn func(Any) (Any, error)) Buildin {
	return func(objs []Any) (Any, error) {
		if len(objs) == 1 {
			return fn(objs[0])
		}

		var out List
		for _, x := range objs {
			result, err := fn(x)
			if err != nil {
				return out, err
			}
			out = append(out, result)
		}
		return out, nil
	}
}
