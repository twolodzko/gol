package parser

type List struct {
	list []interface{}
}

func (l *List) Push(obj interface{}) {
	l.list = append(l.list, obj)
}

func NewList(objs ...interface{}) List {
	var list List
	for _, obj := range objs {
		list.Push(obj)
	}
	return list
}

func Parse(str string) (List, error) {
	var list List

	// for i, ch := range str {

	// }

	return list, nil
}
