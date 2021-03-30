package evaluator

import (
	"github.com/twolodzko/goal/objects"
)

var baseEnv = map[string]objects.Object{
	"true":  True,
	"false": False,
}
