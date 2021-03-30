package evaluator

import "github.com/twolodzko/goal/objects"

type Buildin = func(objects.List) (objects.Object, error)

var buildins = map[string]Buildin{
	"str": str,
}

func str(o objects.List) (objects.Object, error) {
	var out []objects.Object
	for _, elem := range o.Val {
		out = append(out, objects.String{Val: elem.String()})
	}
	return objects.List{Val: out}, nil
}
